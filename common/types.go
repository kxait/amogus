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
