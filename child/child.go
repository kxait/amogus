package child

import (
	"amogus/pvm"
	"amogus/pvm_rpc"
	"fmt"
	"strconv"
	"sync"
	"time"
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

func RunChild() error {
	parentId, err := pvm.Parent()

	if err != nil {
		panic(err)
	}

	srv := &pvm_rpc.RpcServer{Handlers: make(map[pvm_rpc.MessageType]pvm_rpc.RpcHandler)}
	srv.Handlers["ping"] = func(m *pvm_rpc.Message) (*pvm_rpc.Message, error) {
		fmt.Printf("request from %d: %s\n", m.CallerTaskId, m.Content)
		return m.CreateResponse("pong"), nil
	}
	srv.Handlers["multiply"] = func(m *pvm_rpc.Message) (*pvm_rpc.Message, error) {
		i1, i2 := 0, 0
		fmt.Sscanf(m.Content, "%d %d", &i1, &i2)

		return m.CreateResponse(strconv.Itoa(i1 * i2)), nil
	}

	fmt.Println("RPC server up and running")

	parent := pvm_rpc.NewTarget(parentId)

	var wg sync.WaitGroup
	wg.Add(2)

	var loopErr error
	go (func() {
		for loopErr == nil {
			time.Sleep(10 * time.Millisecond)
			loopErr = srv.StepEventLoop()
		}
		wg.Done()
	})()

	var clientErr error
	go (func() {
		for clientErr == nil {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("calling multiply...\n")
			res := <-parent.Call("multiply", "21 37")
			clientErr = res.Err
			fmt.Printf("multiply response: %+v\n", res.Response)
		}
		wg.Done()
	})()

	wg.Wait()

	fmt.Printf("client stopped: %s or %s\n", loopErr.Error(), clientErr.Error())

	return nil
}
