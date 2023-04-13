package cracker

import (
	"amogus/common"
	"sync"
)

type safeCounter struct {
	mut  sync.Mutex
	list []common.HashPair
}

type safeCounterMap struct {
	mut sync.Mutex
	ma  *map[string]string
}
