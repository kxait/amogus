package child

type ChildState int

const (
	Start          ChildState = iota + 1
	ConfigReceived ChildState = iota + 2
	HashesReceived ChildState = iota + 3
	Idle           ChildState = iota + 4
	Cracking       ChildState = iota + 5
)
