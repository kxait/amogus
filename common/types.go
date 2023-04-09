package common

import "amogus/pvm_rpc"

type ChildState int

const (
	Start          ChildState = iota + 1
	ConfigReceived ChildState = iota + 2
	HashesReceived ChildState = iota + 3
	Idle           ChildState = iota + 4
	Cracking       ChildState = iota + 5
)

const (
	GetConfig         pvm_rpc.MessageType = "getConfig"
	GetHashesInfo     pvm_rpc.MessageType = "getHashesInfo"
	GetHashesPart     pvm_rpc.MessageType = "getHashesPart"
	HashCracked       pvm_rpc.MessageType = "hashCracked"
	GetNextAssignment pvm_rpc.MessageType = "getNextAssignment"
)
