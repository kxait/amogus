package main

import (
	"amogus/pvm"
	"fmt"
	"os"
	"runtime"
)

func TestPvm() {
	defer pvm.Exit()

	hostname, err := os.Hostname()
	exitIfError(err)

	myId, err := pvm.Mytid()
	exitIfError(err)

	fmt.Printf("[%x] startuje\n", myId)

	parentId, err := pvm.Parent()

	if isPvmNoParent(err) {
		err := pvm.CatchoutStdout()
		exitIfError(err)

		fmt.Printf("[%x] tworze dzieci\n", myId)

		spawnResult, err := pvm.Spawn("amogus", []string{"--test"}, pvm.TaskDefault, "", 3)
		if err != nil {
			fmt.Printf("[%x] nie moge stworzyc dzieci\n", myId)
			err := pvm.Perror("pvm_spawn")
			exitIfError(err)
			os.Exit(1)
		}

		fmt.Printf("[%x] stworzylem pomyslnie %d procesy\n", myId, spawnResult.Numt)

		for i := 0; i < spawnResult.Numt; i++ {
			fmt.Printf("Dziecko nr %d, tId: %x\n", i, spawnResult.TIds[i])
			_, err := pvm.Initsend(pvm.DataDefault)

			exitIfError(err)
			_, err = pvm.PackfString("%s", "Czesc, dzialasz ?")
			exitIfError(err)

			err = pvm.Send(spawnResult.TIds[i], i+1)
			exitIfError(err)
		}
		for i := 0; i < spawnResult.Numt; i++ {
			pvm.Recv(spawnResult.TIds[i], -1)
			unpackf_result, err := pvm.UnpackfString("%s", 1024)
			exitIfError(err)

			fmt.Printf("Dziecko wysyla: %s\n", unpackf_result)
		}
	} else if parentId > 0 {
		fmt.Printf("[%x] Dziecko: moj rodzic %x\n", myId, parentId)

		bufId, err := pvm.Recv(parentId, -1)
		exitIfError(err)

		bufinfoResult, err := pvm.Bufinfo(bufId)
		exitIfError(err)

		fmt.Printf("Wiadomosc ma %d bajty, id: %x; tid nadawcy: %x\n",
			bufinfoResult.Bytes,
			bufinfoResult.MsgTag,
			bufinfoResult.TId)

		unpackfResult, err := pvm.UnpackfString("%s", 1024)
		exitIfError(err)

		fmt.Printf("%s\n", unpackfResult)

		_, err = pvm.Initsend(pvm.DataDefault)
		exitIfError(err)

		_, err = pvm.PackfString("%s", hostname)
		exitIfError(err)

		err = pvm.Send(parentId, bufinfoResult.TId)
		exitIfError(err)

		fmt.Printf("Test")
	} else {
		exitIfError(err)
	}
}

func exitIfError(err error) {
	if err != nil {
		_, _, no, ok := runtime.Caller(1)
		if ok {
			fmt.Printf("error from test.go#%d\n", no)
		}
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func isPvmNoParent(err error) bool {
	switch e := err.(type) {
	case *pvm.PvmError:
		return e.ErrorCode == pvm.NoParent
	}
	return false
}
