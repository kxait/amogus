package child

import (
	"amogus/common"
	"amogus/config"
	"amogus/parent"
	"amogus/pvm"
	"amogus/pvm_rpc"
	"amogus/tstree"
	"bufio"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

/*
child states:
  - start
  - config received
  - hashes received
  - idle
    ^ |
    | v
  - cracking

steps for changing state:
  - start: initial
  - config received: parent sent over the config parameters
  - hashes received: parent sent over the entire hashes file
  - cracking: parent sent over the last hashed origin and cracking has begun
  - idle: cracking has finished (chunk was finished) and awaiting next step from parent
*/

var hashesPath string = "/tmp/hashes_to_crack"

type childState struct {
	currentAssignment string
	currentState      common.ChildState
	config            config.AmogusConfig
	hashesInfo        config.HashesInfo
	hashPartReceived  int64
}

func RunChild() error {
	state := childState{
		currentState: common.Start,
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
			//time.Sleep(1 * time.Millisecond)
			clientErr = work(&state, parent)
		}
		wg.Done()
	})()

	wg.Wait()

	fmt.Printf("client stopped: %s\n", clientErr.Error())

	return nil
}

func work(state *childState, parent *pvm_rpc.Target) error {
	if state.currentState == common.Start {
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

		state.currentState = common.ConfigReceived
		state.hashPartReceived = -1
	} else if state.currentState == common.ConfigReceived {
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
			state.currentState = common.HashesReceived
			return nil
		}
	} else if state.currentState == common.HashesReceived {
		state.currentState = common.Idle
	} else if state.currentState == common.Cracking {
		hashes := generateHashes(&state.config, state.currentAssignment, state.config.ChunkSize)
		cracked := findStringsInFile(hashesPath, hashes)

		for _, c := range cracked {
			res := <-parent.Call(common.HashCracked, fmt.Sprintf("%s %s", c.hash, c.origin))
			if res.Err != nil {
				return res.Err
			}
		}

		state.currentState = common.Idle
	} else if state.currentState == common.Idle {
		res := <-parent.Call(common.GetNextAssignment, "")
		if res.Err != nil {
			return res.Err
		}

		state.currentAssignment = res.Response.Content
		fmt.Printf("got assignment: %s chunk size %d\n", state.currentAssignment, state.config.ChunkSize)
		state.currentState = common.Cracking
	}
	return nil
}

func findStringsInFile(filename string, sussy []hashPair) []hashPair {
	var result []hashPair
	file, _ := os.Open(filename)
	defer file.Close()

	lut := tstree.BuildLookupTableFromLines([]string{})
	m := make(map[string]string, 0)
	for _, sus := range sussy {
		//fmt.Printf("%+v\n", sus)
		m[sus.hash] = sus.origin
		lut.AppendLine(sus.hash)
	}

	counter := safeCounter{list: make([]hashPair, 0)}
	var wg sync.WaitGroup

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		go (func(l *tstree.LookupTable, ll string, c *safeCounter, ma *map[string]string) {
			if l.Has(ll) {
				//fmt.Printf("FOUND! hash %s origin %s", line, m[line])
				c.mut.Lock()
				c.list = append(c.list, hashPair{hash: ll, origin: (*ma)[ll]})
				c.mut.Unlock()
			}
			wg.Done()
		})(lut, line, &counter, &m)
	}

	wg.Wait()

	return result
}

type hashPair struct {
	hash   string
	origin string
}

type safeCounter struct {
	mut  sync.Mutex
	list []hashPair
}

func generateHashes(cfg *config.AmogusConfig, last string, amount int) []hashPair {
	lastLast := last
	var origins []string
	for i := 0; i < amount; i++ {
		lastLast = parent.GetNextValue(cfg, lastLast)
		origins = append(origins, lastLast)
	}

	c := safeCounter{list: make([]hashPair, 0)}

	var wg sync.WaitGroup

	for _, origin := range origins {
		wg.Add(1)
		go (func(origin string, c *safeCounter) {
			if cfg.Mode != config.Sha512 {
				panic("at the disco")
			}
			hash := hashSha512(origin)
			c.mut.Lock()
			c.list = append(c.list, *hash)
			c.mut.Unlock()
			wg.Done()
		})(origin, &c)
	}

	wg.Wait()

	return c.list
}

func hashSha512(origin string) *hashPair {
	bytes := []byte(origin)
	sha := sha512.New()
	sha.Write(bytes)

	hash := sha.Sum(nil)
	result := &hashPair{
		hash:   fmt.Sprintf("%x", hash[:]),
		origin: origin,
	}

	return result
}
