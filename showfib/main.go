package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/samonzeweb/profilinggo/fibonacci"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage : showfib n")
		os.Exit(1)
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// In real life, don't write to os.Stderr ;)
	// CPU profiling
	pprof.StartCPUProfile(os.Stderr)
	defer pprof.StopCPUProfile()

	fmt.Println(fibonacci.Suite(n))
}
