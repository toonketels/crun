crun - continuously Compile and Run
===================================

Crun is a command line tool that compiles and runs your go code and watches all .go files in the current directory,
recompiling and rerunning your code when the files change.

Useful to start webapps and recompiling them each time a file changes. As such it is an alternative to 
`go run server.go`, each time terminating the command yourself and doing `go run server.go` again.

[![Build Status](https://drone.io/github.com/toonketels/crun/status.png)](https://drone.io/github.com/toonketels/crun/latest)


Install
---------------------------


Ensure Go is installed and your [GOPATH](http://golang.org/doc/code.html#GOPATH) set.

Install package:

    go get github.com/toonketels/crun



Usage
----------------------------------

	// Change to directory with your app go source code
	cd /path/to/my/app

    // Start compiling and running app
    crun


You can pass additional arguments to the binary.

    // additional arguments after `--` are passed to the binary
	crun -- --port=:3000