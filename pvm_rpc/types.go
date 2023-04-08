package pvm_rpc

type Target struct {
	TaskId int
}

func NewTarget(taskId int) *Target {
	targ := Target{TaskId: taskId}
	//targ.ResetSendBuffer()

	return &targ
}

type ReceiveResult struct {
	Response *MessageResponse
	Err      error
}

type MessageType string
type Message struct {
	// IMPORTANT: Generated randomly by the caller
	Id int
	// Will not change on response (will be client tId)
	CallerTaskId int
	Type         MessageType
	Content      string
}

func (m *Message) CreateResponse(content string) *Message {
	response := *m
	response.Content = content
	return &response
}

type MessageResponse struct {
	Message
	IsError bool
}

type RpcHandler func(*Message) (*Message, error)

type RpcServer struct {
	Handlers map[MessageType]RpcHandler
}
