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

type GetNextAssignmentArgs struct {
	ChunkTimeMillis int64
}

type HashCrackedArgs struct {
	Hash   string
	Origin string
}

type GetHashesPartArgs struct {
	Part int
}

type ShadowMode int

const (
	ShadowSha512 ShadowMode = iota
)

type HashPair struct {
	Hash   string
	Origin string
}

type ChildState int

const (
	Start ChildState = iota
	ConfigReceived
	HashesReceived
	Idle
	Cracking
)
