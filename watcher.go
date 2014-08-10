package main

import (
	"Path"
	"code.google.com/p/go.exp/fsnotify"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)

// Watch watches in the current directory for file changes
// and notifies the command dispatcher.
func watch(dispatchChan chan string) {

	// Start a file system watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher", err)
	}

	signalSourceChanged := debounce(800*time.Millisecond, func() {
		dispatchChan <- SOURCECHANGED
	})

	// Start waiting for file system messages to receive...
	go func() {
		for {
			select {
			// On any event
			case ev := <-watcher.Event:
				if isGoFile(ev.Name) && (ev.IsCreate() || ev.IsDelete()) {
					signalSourceChanged()
				}

			// Stop when encountering errors...
			case err := <-watcher.Error:
				log.Fatal("Error start watching", err)
			}
		}
	}()

	// Create a list of all directories to watch...
	directories := append(findDirectoriesIn(DIR_TO_WATCH), DIR_TO_WATCH)
	for _, directory := range directories {
		// Configure to watcher to watch the files we want to...
		if err := watcher.Watch(directory); err != nil {
			log.Fatal("Error start watching", err)
		}
	}
}

func isGoFile(path string) bool {
	regex := regexp.MustCompile(`.*\.go$`)
	return regex.MatchString(path)
}

// Recursively finds directories within the given directory
func findDirectoriesIn(dir string) (dirs []string) {

	// Get all files in current dir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("Error reading dir", err)
	}

	// Investigate each file...
	for _, file := range files {

		// Create the path starting from the starting dir (current directory)
		path := path.Join(dir, file.Name())

		// We want non hidden directories (no .git)...
		if file.IsDir() && !strings.HasPrefix(file.Name(), DIR_TO_WATCH) {

			// Aggregate them and go deeper...
			dirs = append(dirs, path)
			dirs = append(dirs, findDirectoriesIn(path)...)
		}
	}
	return
}
