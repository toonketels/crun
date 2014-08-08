package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Task struct {
	Name      string
	cmd       *exec.Cmd
	done      chan bool
	do        func()
	IsRunning bool
	proc      *os.Process
}

func (task *Task) Wait() {
	_ = <-task.done
	return
}

func (task *Task) Start() *Task {
	task.done = make(chan bool, 1)
	go func() {
		task.IsRunning = true
		log.Println("TASK", task.Name, ": started")
		task.do()
		task.IsRunning = false
		task.done <- true
		log.Println("TASK", task.Name, ": finished")

	}()
	return task
}

func (task *Task) Kill() *Task {

	if err := task.proc.Kill(); err != nil {
		log.Println("Killing process erring", err)
	}
	_, err := task.proc.Wait()
	if err != nil {
		log.Println("Wait process erring")
	}
	task.proc = nil
	task.Wait()
	return task
}

// Compile creates and returns a compile task
func createCompileTask() (task *Task) {

	task = new(Task)
	task.Name = "COMPILE"
	task.IsRunning = false
	task.do = task.compile

	return
}

func createRunTask() (task *Task) {

	task = new(Task)
	task.Name = "RUN"
	task.IsRunning = false
	task.do = task.run

	return
}

func (task *Task) run() {

	// Build command the execute.
	cmd := exec.Command(BIN, flag.Args()...)

	// Create a pipe from cmd stdarr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("STDERR", err)
	}

	// Create a pipe from cmd stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("STDOUT", err)
	}

	// Start the cmd
	err = cmd.Start()
	if err != nil {
		log.Fatal("ERR", err)
	}
	task.proc = cmd.Process

	// Pipe the cmd stdout into our program's stdout,
	// making it visible to us
	io.Copy(os.Stdout, stdout)

	// Read all reported erros from the cmd
	errBuf, _ := ioutil.ReadAll(stderr)

	// Wait for the cmd to done
	_ = cmd.Wait()

	// If we had some errors...
	// We don't check cmd.Wait errors because it will always output
	// an error when we kill the process
	if len(errBuf) != 0 {
		// Dump all its stdout output
		log.Println("ERRBUF", string(errBuf))
	}
}

func (task *Task) compile() {
	// Build command the execute.
	cmd := exec.Command("go", "build", "-o", "CRUN_BIN.tmp", ".")

	// Create a pipe from cmd stdarr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal("STDERR", err)
	}

	// Create a pipe from cmd stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("STDOUT", err)
	}

	// Start the cmd
	err = cmd.Start()
	if err != nil {
		log.Fatal("ERR", err)
	}
	task.proc = cmd.Process

	// Pipe the cmd stdout into our program's stdout,
	// making it visible to us
	io.Copy(os.Stdout, stdout)

	// Read all reported erros from the cmd
	errBuf, _ := ioutil.ReadAll(stderr)

	// Wait for the cmd to done
	err = cmd.Wait()

	// We don't check cmd.Wait errors because it will always output
	// an error when we kill the process
	if len(errBuf) != 0 {
		// Dump all its stdout output
		log.Fatal("ERRBUF", string(errBuf))
	}
}
