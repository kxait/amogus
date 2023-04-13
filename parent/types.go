package parent

import "sync"

const hashRatesCount int = 10

type hashRateSet struct {
	hashRatesByTid map[int]int64
	pushMutex      sync.Mutex
}

func (h *hashRateSet) pushHashRate(tid int, hashRate int64) {
	h.pushMutex.Lock()
	h.hashRatesByTid[tid] = hashRate
	h.pushMutex.Unlock()
}

func (h *hashRateSet) getHashRate() int64 {
	var sum int64
	for _, hr := range h.hashRatesByTid {
		if hr != 0 {
			sum += hr
		}
	}

	return sum
}

func (h *hashRateSet) init() {
	h.hashRatesByTid = make(map[int]int64)
	h.pushMutex = sync.Mutex{}
}
