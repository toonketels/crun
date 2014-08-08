crun - continuously Compile and Run
===================================


The `crun` command compiles and runs your code and watches all .go files in the current directory,
recompiling and rerunning your code when the files change.


Usage
----------------------------------

    // compiles the source in the current dir and executes the binary
    crun

    // additional arguments are passed to the binary
	crun server.go --port=:3000