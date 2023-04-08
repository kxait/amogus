package pvm_rpc

import "sync"

var instance = sync.Mutex{}
var once sync.Once

func GetMutex() *sync.Mutex {
	once.Do(func() {
		instance = sync.Mutex{}
	})

	return &instance
}
