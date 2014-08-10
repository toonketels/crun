package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

const (
	START         = "start"
	SOURCECHANGED = "sourcechanged"

	DIR_TO_WATCH = "."
	BIN          = "./CRUN_BIN.tmp"
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

		dispatchChan := createDispatcher()
		watch(dispatchChan)
		dispatchChan <- START

		waitToTerminate()
	}
}

func createDispatcher() (dispatchChan chan string) {

	dispatchChan = make(chan string)

	go func() {
		compile := createCompileTask()
		run := createRunTask()

		for {
			msg := <-dispatchChan
			switch msg {
			case START:

				compile.Start()
				compile.Wait()
				run.Start()

			case SOURCECHANGED:

				compile.Start()
				compile.Wait()
				if run.IsRunning {
					run.Kill()
				}
				run.Start()
			}
		}
	}()

	return
}

func waitToTerminate() {
	// / Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block for signals to receive
	_ = <-c
	fmt.Println("")
	log.Println("CRUN EXITING")
	_ = os.Remove(BIN)
}
