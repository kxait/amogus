package cracker

import "sync"

type HashPair struct {
	Hash   string
	Origin string
}

type safeCounter struct {
	mut  sync.Mutex
	list []HashPair
}

type safeCounterMap struct {
	mut sync.Mutex
	ma  *map[string]string
}
