package parent

import (
	"amogus/common"
	"amogus/config"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"time"

	pvm_rpc "github.com/kxait/pvm-rpc"
	"github.com/kxait/pvm-rpc/pvm"
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
	ranOut     bool
	shadowMode common.ShadowMode

	hashrate *hashRateSet
}

func RunParent(hashesPath string, configPath string, output string) error {

	state := parentState{lastOrigin: "", hashrate: &hashRateSet{}}
	state.hashrate.init()

	fmt.Println(amogus)

	fmt.Printf("hashes: %s config: %s output: %s\n", hashesPath, configPath, output)

	cfg, err := config.GetConfig(configPath)
	if err != nil {
		return err
	}

	err, shadowMode := config.ValidateHashes(hashesPath, cfg.Mode)
	if err != nil {
		return err
	}

	if shadowMode != nil {
		state.shadowMode = *shadowMode
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
		if pvmErr, ok := err.(*pvm.PvmError); ok {
			if(pvmErr.ErrorCode > 0) {
				fmt.Printf("some tasks failed to spawn: %+v\n", res.TIds)
			}
		}
		debug.PrintStack()
		return err
	}

	defer (func() {
		for _, c := range res.TIds {
			pvm.Kill(c)
		}
	})()

	var wg sync.WaitGroup
	wg.Add(1)

	var loopErr error
	go (func() {
		for loopErr == nil {
			time.Sleep(10 * time.Millisecond)
			loopErr = srv.StepEventLoop()

			// if state.ranOut == true {
			// 	loopErr = fmt.Errorf("finished!")
			// 	break
			// }
		}
		debug.PrintStack()
		wg.Done()
	})()

	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan

		die(res)
	}()

	go (func() {
		oa, err := config.CreateOutputAppender("hashrate")
		if err != nil {
			panic(err)
		}
		i := 0
		for {
			time.Sleep(5 * time.Second)
			hashrate := state.hashrate.getHashRate()
			fmt.Printf("[HASHRATE] %d h/s (%+v)\n", state.hashrate.getHashRate(), state.hashrate.hashRatesByTid)
			oa(fmt.Sprintf("%d %d %+v", i, hashrate, state.hashrate.hashRatesByTid))
			i++
			if cfg.TestSuiteSampleSize > 0 && i >= cfg.TestSuiteSampleSize {
				die(res)
			}
		}
	})()

	wg.Wait()

	fmt.Printf("server stopped: %s\n", loopErr.Error())

	return nil
}

func die(spawnResult *pvm.Spawn_result) {
	for _, c := range spawnResult.TIds {
		pvm.Kill(c)
	}

	os.Exit(0)
}
