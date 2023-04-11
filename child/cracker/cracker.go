package cracker

import (
	"amogus/common"
	"amogus/config"
	"bufio"
	"os"
	"sync"
)

func FindStringsInFile(filename string, sussy *map[string]string) []HashPair {
	var result []HashPair
	file, _ := os.Open(filename)
	defer file.Close()

	counter := safeCounter{list: make([]HashPair, 0)}
	var wg sync.WaitGroup

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		go func(ll string, c *safeCounter, ma *map[string]string) {
			if origin, ok := (*ma)[ll]; ok {
				c.mut.Lock()
				c.list = append(c.list, HashPair{Hash: ll, Origin: origin})
				c.mut.Unlock()
			}
			wg.Done()
		}(line, &counter, sussy)
	}

	wg.Wait()

	return result
}

func GenerateHashes(cfg *config.AmogusConfig, last string, amount int) *map[string]string {
	lastLast := last
	var origins []string
	for i := 0; i < amount; i++ {
		lastLast = common.GetNextValue(cfg, lastLast)
		origins = append(origins, lastLast)
	}

	result := make(map[string]string)

	c := safeCounterMap{ma: &result}

	var wg sync.WaitGroup

	for _, origin := range origins {
		wg.Add(1)
		go (func(origin string, c *safeCounterMap) {
			if cfg.Mode != config.Sha512 {
				panic("at the disco")
			}
			hash := hashSha512(origin)
			c.mut.Lock()
			(*c.ma)[hash.Hash] = origin
			c.mut.Unlock()
			wg.Done()
		})(origin, &c)
	}

	wg.Wait()

	return c.ma
}
