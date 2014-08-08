package main

import (
	"time"
)

// Debounce creates and returns a new debounced version of the passed function
// which will postpone its execution until after wait Duration have elapsed
// since the last time it was invoked.
func debounce(wait time.Duration, fn func()) func() {
	var timer *time.Timer

	// Wrap the the function so we can discard the timer when it's done.
	fnToCall := func() {
		timer = nil
		fn()
	}

	// Return a function that can be called multiple times
	return func() {
		// Only schedule execution of the function when no timer is present.
		if timer != nil {
			return
		}
		// Schedules the execution of the function.
		timer = time.AfterFunc(wait, fnToCall)
	}
}
