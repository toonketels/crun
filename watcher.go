package main

import (
	"code.google.com/p/go.exp/fsnotify"
	"log"
	"regexp"
	"time"
)

func watch(msgBus chan string) {

	// Start a file system watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher", err)
	}

	signalSourceChanged := debounce(2*time.Second, func() {
		msgBus <- SOURCECHANGED
	})

	// Start waiting for file sytem messages to recieve...
	go func() {
		for {
			select {
			// On any event
			case ev := <-watcher.Event:
				if isGoFile(ev.Name) && (ev.IsCreate() || ev.IsDelete()) {
					signalSourceChanged()
				}

			// Stop when encoutering errors...
			case err := <-watcher.Error:
				log.Fatal("Error start watching", err)
			}
		}
	}()

	// Configure to watcher to watch the files we want to...
	if err := watcher.Watch("."); err != nil {
		log.Fatal("Error start watching", err)
	}
	// Configure to watcher to watch the files we want to...
	if err := watcher.Watch("structs"); err != nil {
		log.Fatal("Error start watching", err)
	}
	// Configure to watcher to watch the files we want to...
	if err := watcher.Watch("lib"); err != nil {
		log.Fatal("Error start watching", err)
	}
}

func isGoFile(path string) bool {
	regex := regexp.MustCompile(`.*\.go$`)
	return regex.MatchString(path)
}
