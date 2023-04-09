package parent

import (
	"amogus/config"
	"amogus/pvm"
	"amogus/pvm_rpc"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

const amogus string = `
           ⣠⣤⣤⣤⣤⣤⣶⣦⣤⣄⡀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⡿⠛⠉⠙⠛⠛⠛⠛⠻⢿⣿⣷⣤⡀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⠋⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⠈⢻⣿⣿⡄⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⣸⣿⡏⠀⠀⠀⣠⣶⣾⣿⣿⣿⠿⠿⠿⢿⣿⣿⣿⣄⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⣿⣿⠁⠀⠀⢰⣿⣿⣯⠁⠀⠀⠀⠀⠀⠀⠀⠈⠙⢿⣷⡄⠀
⠀⠀⣀⣤⣴⣶⣶⣿⡟⠀⠀⠀⢸⣿⣿⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣷⠀
⠀⢰⣿⡟⠋⠉⣹⣿⡇⠀⠀⠀⠘⣿⣿⣿⣿⣷⣦⣤⣤⣤⣶⣶⣶⣶⣿⣿⣿⠀
⠀⢸⣿⡇⠀⠀⣿⣿⡇⠀⠀⠀⠀⠹⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠃⠀
⠀⣸⣿⡇⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠉⠻⠿⣿⣿⣿⣿⡿⠿⠿⠛⢻⣿⡇⠀⠀
⠀⣿⣿⠁⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣧⠀⠀
⠀⣿⣿⠀⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⠀⠀
⠀⣿⣿⠀⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⠀⠀
⠀⢿⣿⡆⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⡇⠀⠀
⠀⠸⣿⣧⡀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⠃⠀⠀
⠀⠀⠛⢿⣿⣿⣿⣿⣇⠀⠀⠀⠀⠀⣰⣿⣿⣷⣶⣶⣶⣶⠶⠀⢠⣿⣿⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⣿⣿⠀⠀⠀⠀⠀⣿⣿⡇⠀⣽⣿⡏⠁⠀⠀⢸⣿⡇⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⣿⣿⠀⠀⠀⠀⠀⣿⣿⡇⠀⢹⣿⡆⠀⠀⠀⣸⣿⠇⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⢿⣿⣦⣄⣀⣠⣴⣿⣿⠁⠀⠈⠻⣿⣿⣿⣿⡿⠏⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠈⠛⠻⠿⠿⠿⠿⠋⠁`

type parentState struct {
	lastOrigin string
}

func RunParent(hashesPath string, configPath string, output string) error {

	state := parentState{}

	fmt.Println(amogus)

	fmt.Printf("hashes: %s config: %s output: %s\n", hashesPath, configPath, output)

	cfg, err := config.GetConfig(configPath)

	if err != nil {
		return err
	}

	oa, err := config.CreateOutputAppender(output)
	if err != nil {
		return err
	}

	srv := &pvm_rpc.RpcServer{Handlers: make(map[pvm_rpc.MessageType]pvm_rpc.RpcHandler)}
	registerParentHandlers(srv, cfg, hashesPath, oa, &state)

	fmt.Println("RPC server up and running")

	pvm.CatchoutStdout()
	res, err := pvm.Spawn("amogus", []string{"--child"}, pvm.TaskDefault, "", int(cfg.Slaves))
	if err != nil {
		return err
	}

	defer (func() {
		for _, c := range res.TIds {
			pvm.Kill(c)
		}
	})()

	//child := pvm_rpc.NewTarget(res.TIds[0])

	var wg sync.WaitGroup
	//wg.Add(2)
	wg.Add(1)

	var loopErr error
	go (func() {
		for loopErr == nil {
			time.Sleep(10 * time.Millisecond)
			loopErr = srv.StepEventLoop()
		}
		wg.Done()
	})()

	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan

		for _, c := range res.TIds {
			pvm.Kill(c)
		}

		// do last actions and wait for all write operations to end

		os.Exit(0)
	}()

	// var clientErr error
	// go (func() {
	// 	for clientErr == nil {
	// 		time.Sleep(10 * time.Millisecond)

	// 	}
	// 	wg.Done()
	// })()

	wg.Wait()

	//fmt.Printf("server stopped: %s or %s\n", loopErr.Error(), clientErr.Error())
	fmt.Printf("server stopped: %s\n", loopErr.Error())

	// last := "aaaaaaaa"
	// var count int64
	// for count < 100000000 {
	// 	//fmt.Println("generating")
	// 	chunk := generateHashes(cfg, last, cfg.ChunkSize)
	// 	count++
	// 	//fmt.Printf("%+v\n", chunk)

	// 	foundHashes := findStringsInFile(hashesPath, chunk)
	// 	for _, foundHash := range foundHashes {
	// 		fmt.Printf("FOUND! hash %s origin %s", foundHash.hash, foundHash.origin)
	// 	}

	// 	last = chunk[len(chunk)-1].origin

	//}

	return nil
}

// func findStringsInFile(filename string, sussy []hashPair) []hashPair {
// 	var result []hashPair
// 	file, _ := os.Open(filename)
// 	defer file.Close()

// 	lut := tstree.BuildLookupTableFromLines([]string{})
// 	m := make(map[string]string, 0)
// 	for _, sus := range sussy {
// 		//fmt.Printf("%+v\n", sus)
// 		m[sus.hash] = sus.origin
// 		lut.AppendLine(sus.hash)
// 	}

// 	counter := safeCounter{list: make([]hashPair, 0)}
// 	var wg sync.WaitGroup

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		wg.Add(1)
// 		go (func(l *tstree.LookupTable, ll string, c *safeCounter, ma *map[string]string) {
// 			if l.Has(ll) {
// 				//fmt.Printf("FOUND! hash %s origin %s", line, m[line])
// 				c.mut.Lock()
// 				c.list = append(c.list, hashPair{hash: ll, origin: (*ma)[ll]})
// 				c.mut.Unlock()
// 			}
// 			wg.Done()
// 		})(lut, line, &counter, &m)
// 	}

// 	wg.Wait()

// 	return result
// }

// type hashPair struct {
// 	hash   string
// 	origin string
// }

// type safeCounter struct {
// 	mut  sync.Mutex
// 	list []hashPair
// }

// func generateHashes(cfg *config.AmogusConfig, last string, amount int) []hashPair {
// 	lastLast := last
// 	var origins []string
// 	for i := 0; i < amount; i++ {
// 		lastLast = GetNextValue(cfg, lastLast)
// 		origins = append(origins, lastLast)
// 	}

// 	c := safeCounter{list: make([]hashPair, 0)}

// 	var wg sync.WaitGroup

// 	for _, origin := range origins {
// 		wg.Add(1)
// 		go (func(origin string, c *safeCounter) {
// 			if cfg.Mode != config.Sha512 {
// 				panic("at the disco")
// 			}
// 			hash := hashSha512(origin)
// 			c.mut.Lock()
// 			c.list = append(c.list, *hash)
// 			c.mut.Unlock()
// 			wg.Done()
// 		})(origin, &c)
// 	}

// 	wg.Wait()

// 	return c.list
// }

// func hashSha512(origin string) *hashPair {
// 	bytes := []byte(origin)
// 	sha := sha512.New()
// 	sha.Write(bytes)

// 	hash := sha.Sum(nil)
// 	result := &hashPair{
// 		hash:   fmt.Sprintf("%x", hash[:]),
// 		origin: origin,
// 	}

// 	return result
// }
