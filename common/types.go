package common

import (
	pvm_rpc "github.com/kxait/pvm-rpc"
)

const (
	GetConfig         pvm_rpc.MessageType = "getConfig"
	GetHashesInfo     pvm_rpc.MessageType = "getHashesInfo"
	GetHashesPart     pvm_rpc.MessageType = "getHashesPart"
	HashCracked       pvm_rpc.MessageType = "hashCracked"
	GetNextAssignment pvm_rpc.MessageType = "getNextAssignment"
)

type ShadowMode int

const (
	ShadowSha512 ShadowMode = iota + 1
)

type HashPair struct {
	Hash   string
	Origin string
}

type ChildState int

const (
	Start          ChildState = iota + 1
	ConfigReceived ChildState = iota + 2
	HashesReceived ChildState = iota + 3
	Idle           ChildState = iota + 4
	Cracking       ChildState = iota + 5
)
