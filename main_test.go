package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

// Test the watcher by verifying if it sends the SOURCECHANGED
// messages on the channel passed.
func TestWatcher(t *testing.T) {
	dir := tempMkdir(t, "crun-watcher")

	if err := os.Chdir(dir); err != nil {
		t.Fatal("Failed to change to temp dir")
	}

	c := make(chan string, 2)
	watch(c)

	notified := 0

	go func() {
		for msg := range c {
			if msg == SOURCECHANGED {
				notified++
			}
		}
	}()

	// Should only do one, since debounced
	ioutil.WriteFile("testfile-1.go", []byte("something"), 0644)
	ioutil.WriteFile("testfile-2.go", []byte("something"), 0644)

	time.Sleep(1 * time.Second)
	if notified != 1 {
		t.Fatal("Watcher should have notified 1 time", notified)
	}

	ioutil.WriteFile("testfile-3.go", []byte("something"), 0644)
	time.Sleep(1 * time.Second)
	if notified != 2 {
		t.Fatal("Watcher should have notified 2 times", notified)
	}

	ioutil.WriteFile("testfile-1.md", []byte("something"), 0644)
	time.Sleep(1 * time.Second)
	if notified != 2 {
		t.Fatal("Watcher should only notify on md files", notified)
	}

	ioutil.WriteFile("testfile-4.go", []byte("something"), 0644)
	ioutil.WriteFile("testfile-5.go", []byte("something"), 0644)

	time.Sleep(1 * time.Second)
	if notified != 3 {
		t.Fatal("Watcher should have notified 3 times", notified)
	}
	close(c)
}

// Test dispatcher by sending it the START and SOURCECHANGED
// messages and see if it compiles and runs the test program.
func TestDispatcher(t *testing.T) {

	dir := tempMkdir(t, "crun-tasks")
	createGoProgram(t, dir)

	if err := os.Chdir(dir); err != nil {
		t.Fatal("Failed to change to temp dir")
	}
	c := createDispatcher()

	if _, err := ioutil.ReadFile(BIN); err.Error() != "open ./CRUN_BIN.tmp: no such file or directory" {
		t.Fatal("There should be no binary yet")
	}

	if _, err := ioutil.ReadFile("STARTED"); err.Error() != "open STARTED: no such file or directory" {
		t.Fatal("There should be no STARTED file yet")
	}

	// START msg should start compiling
	c <- START

	time.Sleep(1 * time.Second)
	if _, err := ioutil.ReadFile(BIN); err != nil {
		t.Fatal("There should be a binary")
	}

	if _, err := ioutil.ReadFile("STARTED"); err != nil {
		t.Fatal("There should be a STARTED file", err)
	}

	// Remove the STARTED file so we can check if recompiling and rerunning worked
	if err := os.Remove("STARTED"); err != nil {
		t.Fatal("Error removing STARTED file")
	}

	// SOURCECHANGED msg should start recompiling and rerunning
	c <- SOURCECHANGED
	time.Sleep(1 * time.Second)
	if _, err := ioutil.ReadFile(BIN); err != nil {
		t.Fatal("There should be a binary")
	}

	if _, err := ioutil.ReadFile("STARTED"); err != nil {
		t.Fatal("There should be a STARTED file", err)
	}
}

// Creates a temporary directory
func tempMkdir(t *testing.T, name string) string {
	dir, err := ioutil.TempDir("", name)
	if err != nil {
		t.Fatalf("failed to create test directory: %s", err)
	}
	return dir
}

// Write a source for dummy go program to dmp dir.
func createGoProgram(t *testing.T, dir string) {
	fileName := dir + "/testprogram.go"
	contents := getGoSourceProgram()
	err := ioutil.WriteFile(fileName, []byte(contents), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %s", err)
	}
}

// Returns the program the be written to disk.
func getGoSourceProgram() string {
	contents := `
package main

import "time"
import "io/ioutil"

func main() {
	ioutil.WriteFile("STARTED", []byte("STARTED\n"), 0644)
	// Ensure it wont exit immediately
	time.Sleep(10 * time.Minute)
	ioutil.WriteFile("FINISHED", []byte("FINISHED\n"), 0644)

}
	`
	return contents
}
