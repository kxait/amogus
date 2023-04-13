package cracker

import (
	"amogus/child/state"
	"amogus/common"
	"amogus/config"
	"amogus/next_value"
	"bufio"
	"os"
	"sync"
)

func FindStringsInFile(filename string, generatedHashes *map[string]string) []common.HashPair {
	file, _ := os.Open(filename)
	defer file.Close()

	counter := safeCounter{list: make([]common.HashPair, 0)}
	var wg sync.WaitGroup

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		wg.Add(1)
		go func(ll string, c *safeCounter, ma *map[string]string) {
			if origin, ok := (*ma)[ll]; ok {
				c.mut.Lock()
				c.list = append(c.list, common.HashPair{Hash: ll, Origin: origin})
				c.mut.Unlock()
			}
			wg.Done()
		}(line, &counter, generatedHashes)
	}

	wg.Wait()

	return counter.list
}

func GenerateHashes(s *state.ChildState) *map[string]string {
	lastLast := s.CurrentAssignment
	var origins []string
	for i := 0; i < s.Config.ChunkSize; i++ {
		lastLast = next_value.GetNextValue(&s.Config, lastLast)
		origins = append(origins, lastLast)
	}

	result := make(map[string]string)

	c := safeCounterMap{ma: &result}

	var wg sync.WaitGroup

	for _, origin := range origins {
		wg.Add(1)
		go (func(origin string, c *safeCounterMap) {
			var hash *common.HashPair
			switch s.Config.Mode {
			case config.Sha512:
				hash = hashSha512(origin)
				break
			case config.Sha256:
				hash = hashSha256(origin)
				break
			case config.Shadow:
				hash = hashShadow(origin, s)
				break
			default:
				panic("at the disco")
			}
			c.mut.Lock()
			(*c.ma)[hash.Hash] = origin
			c.mut.Unlock()
			wg.Done()
		})(origin, &c)
	}

	wg.Wait()

	return c.ma
}
