package parent

import "sync"

const hashRatesCount int = 10

type hashRateSet struct {
	hashRates []int64
	pointer   int
	pushMutex sync.Mutex
}

func (h *hashRateSet) pushHashRate(hashRate int64) {
	h.pushMutex.Lock()
	h.hashRates[h.pointer] = hashRate
	h.pointer++
	if h.pointer >= hashRatesCount {
		h.pointer = 0
	}
	h.pushMutex.Unlock()
}

func (h *hashRateSet) getHashRate() int64 {
	var sum int64
	for _, hr := range h.hashRates {
		if hr != 0 {
			sum += hr
		}
	}

	return sum
}

func (h *hashRateSet) init() {
	h.hashRates = make([]int64, hashRatesCount)
	h.pointer = 0
	h.pushMutex = sync.Mutex{}
}
