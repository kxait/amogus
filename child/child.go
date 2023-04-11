package child

import (
	"amogus/child/cracker"
	"amogus/common"
	"amogus/config"
	"amogus/pvm"
	"amogus/pvm_rpc"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var hashesPath string = "/tmp/hashes_to_crack"

type childState struct {
	currentAssignment string
	currentState      ChildState
	config            config.AmogusConfig
	hashesInfo        config.HashesInfo
	hashPartReceived  int64
}

func RunChild() error {
	state := childState{
		currentState: Start,
	}

	parentId, err := pvm.Parent()

	if err != nil {
		panic(err)
	}
	parent := pvm_rpc.NewTarget(parentId)

	me, err := pvm.Mytid()

	if err != nil {
		panic(err)
	}

	hashesPath = fmt.Sprintf("%s+%d", hashesPath, me)

	var wg sync.WaitGroup
	wg.Add(1)

	var clientErr error
	go (func() {
		for clientErr == nil {
			time.Sleep(1 * time.Millisecond)
			clientErr = work(&state, parent)
		}
		wg.Done()
	})()

	wg.Wait()

	fmt.Printf("client stopped: %s\n", clientErr.Error())

	return nil
}

func work(state *childState, parent *pvm_rpc.Target) error {
	if state.currentState == Start {
		res := <-parent.Call(common.GetConfig, "")
		if res.Err != nil {
			return res.Err
		}

		err := json.Unmarshal([]byte(res.Response.Content), &state.config)
		if err != nil {
			return err
		}

		res = <-parent.Call(common.GetHashesInfo, "")
		if res.Err != nil {
			return res.Err
		}

		err = json.Unmarshal([]byte(res.Response.Content), &state.hashesInfo)
		if err != nil {
			return err
		}

		fmt.Println(state.config)
		fmt.Println(state.hashesInfo)

		_, err = config.CreateOutputAppenderWithTruncate(hashesPath)
		if err != nil {
			return err
		}

		state.currentState = ConfigReceived
		state.hashPartReceived = -1
	} else if state.currentState == ConfigReceived {
		res := <-parent.Call(common.GetHashesPart, strconv.FormatInt(state.hashPartReceived+1, 10))
		if res.Err != nil {
			return res.Err
		}

		oa, err := config.CreateOutputAppender(hashesPath)
		if err != nil {
			return err
		}

		oa(res.Response.Message.Content)

		fmt.Printf("%d/%d: len %d\n", state.hashPartReceived+1, state.hashesInfo.Parts, len(res.Response.Message.Content))

		state.hashPartReceived++
		if state.hashPartReceived >= state.hashesInfo.Parts-1 {
			state.currentState = HashesReceived
			return nil
		}
	} else if state.currentState == HashesReceived {
		state.currentState = Idle
	} else if state.currentState == Cracking {
		hashes := cracker.GenerateHashes(&state.config, state.currentAssignment, state.config.ChunkSize)
		cracked := cracker.FindStringsInFile(hashesPath, hashes)

		for _, c := range cracked {
			res := <-parent.Call(common.HashCracked, fmt.Sprintf("%s %s", c.Hash, c.Origin))
			if res.Err != nil {
				return res.Err
			}
		}

		state.currentState = Idle
	} else if state.currentState == Idle {
		res := <-parent.Call(common.GetNextAssignment, "")
		if res.Err != nil {
			return res.Err
		}

		state.currentAssignment = res.Response.Content
		fmt.Printf("got assignment: %s chunk size %d\n", state.currentAssignment, state.config.ChunkSize)
		state.currentState = Cracking
	}
	return nil
}
