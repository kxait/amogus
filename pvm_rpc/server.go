package pvm_rpc

import (
	"amogus/pvm"
	"encoding/json"
	"fmt"
)

func (rs *RpcServer) StepEventLoop() error {
	mutex := GetMutex()

	mutex.Lock()
	msg, err := rs.pullMessage()
	mutex.Unlock()

	if err != nil {
		return err
	}

	if msg == nil {
		return nil
	}

	handler, ok := rs.Handlers[msg.Type]

	if !ok {
		return fmt.Errorf("handler not found for message type %s (id %d)", msg.Type, msg.Id)
	}

	result, err := handler(msg)

	mutex.Lock()
	defer mutex.Unlock()
	if err != nil {
		rs.send(&MessageResponse{
			Message: *msg.CreateResponse(err.Error()),
			IsError: true,
		})
		return nil
	}

	rs.send(&MessageResponse{
		Message: *result,
		IsError: false,
	})

	return nil
}

func (rs *RpcServer) send(msg *MessageResponse) error {
	serialized, err := json.Marshal(*msg)

	if err != nil {
		return err
	}

	pvm.Initsend(pvm.DataDefault)

	pvm.PackfString("%s", string(serialized))

	pvm.Send(msg.CallerTaskId, msg.Id)

	return nil
}

func (rs *RpcServer) pullMessage() (*Message, error) {
	bufId, err := pvm.Nrecv(-1, 0)

	if err != nil {
		return nil, err
	}

	// did not receive anything
	if bufId == 0 {
		return nil, nil
	}

	data, err := pvm.UnpackfString("%s", MaxPacketSize)

	if err != nil {
		return nil, err
	}

	msg := &Message{}

	err = json.Unmarshal([]byte(data), &msg)

	if err != nil {
		return nil, err
	}

	return msg, nil
}
