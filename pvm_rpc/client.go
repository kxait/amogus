package pvm_rpc

import (
	"amogus/pvm"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"
)

func Sus() {}

func (t *Target) ResetSendBuffer() {
	pvm.Initsend(pvm.DataDefault)
}

func (t *Target) Kill() error {
	return pvm.Kill(t.TaskId)
}

func (t *Target) Call(msgType MessageType, content string) <-chan *ReceiveResult {
	r := make(chan *ReceiveResult)

	mutex := GetMutex()

	go (func() {
		me, err := pvm.Mytid()

		if err != nil {
			r <- &ReceiveResult{Err: err}
			return
		}

		id := rand.Intn(math.MaxInt32)

		msg := &Message{
			Id:           id,
			CallerTaskId: me,
			Type:         msgType,
			Content:      content,
		}

		mutex.Lock()
		err = t.send(msg)
		mutex.Unlock()
		if err != nil {
			r <- &ReceiveResult{Err: err}
			return
		}

		result := <-t.receiveContinuously(msg.Id)

		if result.Err != nil {
			r <- &ReceiveResult{Err: err}
			return
		}

		r <- result
	})()

	return r
}

func (t *Target) send(msg *Message) error {
	serialized, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	t.ResetSendBuffer()

	_, err = pvm.PackfString("%s", string(serialized))
	if err != nil {
		return err
	}
	err = pvm.Send(t.TaskId, 0)
	if err != nil {
		return err
	}

	return nil
}

func (t *Target) receiveContinuously(id int) <-chan *ReceiveResult {
	r := make(chan *ReceiveResult)
	mutex := GetMutex()

	fin := time.Now().Add(time.Duration(RequestTimeoutSeconds * int(time.Second)))

	go (func() {
		for {
			mutex.Lock()
			bufId, err := pvm.Nrecv(t.TaskId, id)
			if err != nil {
				r <- &ReceiveResult{Err: err}
				mutex.Unlock()
				break
			}

			if bufId == 0 {
				mutex.Unlock()
				if time.Now().After(fin) {
					r <- &ReceiveResult{Err: fmt.Errorf("request timeout")}
					break
				}
				time.Sleep(10 * time.Millisecond)
				continue
			}

			content, err := pvm.UnpackfString("%s", MaxPacketSize)
			mutex.Unlock()

			if err != nil {
				r <- &ReceiveResult{Err: err}
				break
			}

			deserialized := &MessageResponse{}

			err = json.Unmarshal([]byte(content), deserialized)

			if err != nil {
				r <- &ReceiveResult{Err: err}
				break
			}

			r <- &ReceiveResult{Response: deserialized}
			break
		}
	})()

	return r
}
