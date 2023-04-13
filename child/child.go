package child

import (
	"amogus/child/cracker"
	"amogus/child/cracker/shadow"
	"amogus/child/state"
	"amogus/common"
	"amogus/config"
	"os"

	"encoding/json"
	"fmt"
	"sync"
	"time"

	pvm_rpc "github.com/kxait/pvm-rpc"
	"github.com/kxait/pvm-rpc/pvm"
)

var hashesPath string = "/tmp/hashes_to_crack"

func RunChild() error {
	state := state.ChildState{
		CurrentState: common.Start,
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

func work(state *state.ChildState, parent *pvm_rpc.Target) error {
	if state.CurrentState == common.Start {
		res := <-parent.Call(common.GetConfig, "")
		if res.Err != nil {
			return res.Err
		}

		err := json.Unmarshal([]byte(res.Response.Content), &state.Config)

		if err != nil {
			return err
		}

		res = <-parent.Call(common.GetHashesInfo, "")
		if res.Err != nil {
			return res.Err
		}

		err = json.Unmarshal([]byte(res.Response.Content), &state.HashesInfo)
		if err != nil {
			return err
		}

		fmt.Println(state.Config)
		fmt.Println(state.HashesInfo)

		_, err = config.CreateOutputAppenderWithTruncate(hashesPath)
		if err != nil {
			return err
		}

		state.CurrentState = common.ConfigReceived
		state.HashPartReceived = -1
	} else if state.CurrentState == common.ConfigReceived {
		res := <-parent.Call(common.GetHashesPart, serialize(common.GetHashesPartArgs{
			Part: int(state.HashPartReceived + 1),
		}))

		if res.Err != nil {
			return res.Err
		}

		oa, err := config.CreateOutputAppender(hashesPath)
		if err != nil {
			return err
		}

		oa(res.Response.Message.Content)

		fmt.Printf("%d/%d: len %d\n", state.HashPartReceived+1, state.HashesInfo.Parts, len(res.Response.Message.Content))

		state.HashPartReceived++
		if state.HashPartReceived >= state.HashesInfo.Parts-1 {
			state.CurrentState = common.HashesReceived
			return nil
		}
	} else if state.CurrentState == common.HashesReceived {
		if state.HashesInfo.ShadowMode != 0 {
			if state.HashesInfo.ShadowMode == common.ShadowSha512 {
				hashes, err := os.ReadFile(hashesPath)
				if err != nil {
					return err
				}

				state.ShadowCrypter = shadow.GetSaltySha512Crypter(string(hashes))
			} else {
				return fmt.Errorf("unsupported shadow mode %d", state.HashesInfo.ShadowMode)
			}
		}

		state.CurrentState = common.Idle
	} else if state.CurrentState == common.Cracking {
		state.LastChunkStart = time.Now()

		hashes := cracker.GenerateHashes(state)
		cracked := cracker.FindStringsInFile(hashesPath, hashes)

		for _, c := range cracked {
			res := <-parent.Call(common.HashCracked, serialize(common.HashCrackedArgs{
				Hash:   c.Hash,
				Origin: c.Origin,
			}))

			if res.Err != nil {
				return res.Err
			}
		}

		state.CurrentState = common.Idle
	} else if state.CurrentState == common.Idle {
		diff := time.Now().Sub(state.LastChunkStart).Milliseconds()

		res := <-parent.Call(common.GetNextAssignment, serialize(common.GetNextAssignmentArgs{
			ChunkTimeMillis: diff,
		}))
		if res.Err != nil {
			return res.Err
		}

		state.CurrentAssignment = res.Response.Content
		fmt.Printf("got assignment: %s chunk size %d\n", state.CurrentAssignment, state.Config.ChunkSize)
		state.CurrentState = common.Cracking

	}
	return nil
}

func serialize(str interface{}) string {
	serialized, err := json.Marshal(str)

	if err != nil {
		panic(err)
	}

	return string(serialized)
}
