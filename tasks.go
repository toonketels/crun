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
	cmd       func() *exec.Cmd
	done      chan bool
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

		task.Println("started")

		cmd := task.cmd()

		// Create a pipe from cmd stdarr
		stderr, err := cmd.StderrPipe()
		if err != nil {
			task.Fatal(err)
		}

		// Create a pipe from cmd stdout
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			task.Fatal(err)
		}

		// Start the cmd
		err = cmd.Start()
		if err != nil {
			task.Fatal(err)
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
			task.Println(string(errBuf))
		}

		task.IsRunning = false

		task.done <- true

		task.Println("finished")

	}()
	return task
}

func (task *Task) Kill() *Task {

	if err := task.proc.Kill(); err != nil {
		task.Println(err)
	}
	_, err := task.proc.Wait()
	if err != nil {
		task.Println(err)
	}
	task.proc = nil
	task.Wait()
	return task
}

func (task *Task) Println(v interface{}) {
	log.Println("TASK", task.Name, ":", v)
}

func (task *Task) Fatal(v interface{}) {
	log.Fatal("TASK", task.Name, ":", v)
}

func createCompileTask() (task *Task) {

	task = new(Task)
	task.Name = "COMPILE"
	task.IsRunning = false
	task.cmd = func() *exec.Cmd { return exec.Command("go", "build", "-o", "CRUN_BIN.tmp", ".") }

	return
}

func createRunTask() (task *Task) {

	task = new(Task)
	task.Name = "RUN"
	task.IsRunning = false
	task.cmd = func() *exec.Cmd { return exec.Command(BIN, flag.Args()...) }

	return
}
