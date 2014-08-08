package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	STOP          = "stop"
	START         = "start"
	SOURCECHANGED = "sourcechanged"
	TERMINATE     = "terminate"

	DIR_TO_WATCH = "."
)

func main() {

	help := flag.Bool("help", false, "Shows this help message")
	flag.Parse()

	if *help {
		fmt.Println("")
		fmt.Println("crun - continuously Compile and Run")
		fmt.Println("")
		fmt.Println("       compiles and runs your code and watches all .go files")
		fmt.Println("       in the current directory, recompiling and rerunning")
		fmt.Println("       your code when the files change.")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("")
		fmt.Println("crun server.go                   compiles the source in the current")
		fmt.Println("                                 dir and executes the binary")

		fmt.Println("crun server.go -- --port=:3000   arguments after `--` are passed to the binary")
		fmt.Println("")
	} else {

		msgBus := make(chan string)
		go createDispatcher(msgBus)

		watch(msgBus)
		msgBus <- START

		waitUntilInstructedToKill(msgBus)
	}
}

func createDispatcher(msgBus chan string) {

	compile := createCompileTask()
	run := createRunTask()
	removeBin := createRemoveBinTask()

	for {
		msg := <-msgBus
		switch msg {
		case START:
			log.Println("dispatcher start")

			compile.Start()
			compile.Wait()
			run.Start()

		case SOURCECHANGED:
			log.Println("dispatcher sourceChanged")

			compile.Start()
			compile.Wait()
			if run.IsRunning {
				run.Kill()
			}
			run.Start()

		case STOP:
			log.Println("dispatcher stop")

			if compile.IsRunning {
				compile.Kill()
			}
			if run.IsRunning {
				run.Kill()
			}
			removeBin.Start()
			removeBin.Wait()

		case TERMINATE:
			log.Println("dispatcher terminate")
		}
	}
}

func waitUntilInstructedToKill(msgBus chan string) {
	// / Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block for signals to receive
	for sig := range c {
		log.Println("Got signal:", sig)
		msgBus <- STOP

		// After 5 sec interrupt will just stop the program
		time.AfterFunc(1*time.Second, func() {
			signal.Stop(c)
			return
		})
	}
}
