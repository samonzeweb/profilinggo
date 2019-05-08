# Basics of benchmarking, profiling and tracing with Go

## Introduction

This documentation gives an overview of possibilities offered by go tooling to measure performance or collect runtime information. It is not a detailed tutorial about benchmarking, profiling or tracing.

This documentation could also act as a reminder.

Most cases could be tested with the simple code accompanying this documentation, so it is easy to start experimenting these tools. It's also designed to be a live demo.

The main subjects here are :

* Benchmarking : focus is on a particular piece of code, allowing measurement of cpu and/or memory information.
* Profiling : aggregated data collected through sampling during program (or test) execution. Profiling has no timeline.
* Tracing : data collected through events occurring during program (or test) execution. Tracing has a timeline.

Profiling and tracing could apply to benchmarks.

## Benchmarking

Benchmarking is done through the Go testing tools. It's rather simple and well documented.

The primary result of a benchmark is, per tested operation :

* the time it takes.
* the amount of memory allocated on the heap.
* the amount of allocations.

Each benchmark could also be a starting point to profiling or tracing operations.

See the code of `fibonacci/fibonacci_test.go` for a minimal example.

### Running benchmarks

Run the tests :

* without any benchmarks : `go test ./fibonacci`
* with benchmarks (time) : `go test ./fibonacci -bench .`
* with benchmarks (time and memory) : `go test ./fibonacci -bench . -benchmem`

The argument following `-bench` is a regular expression. All banchmark functions whose names match are executed. The `.` in the previous examples isn't the currect directory but a pattern matching all tests. To run a specific benchmark, use the regexp : `-bench Suite` (means *everything containing Suite*).

Useful tip : see `ResetTimer()` to ignore test setup in measures, see also `StopTimer()` and `StartTimer()`: https://golang.org/pkg/testing/#B.ResetTimer

### Comparing benchmarks

It's possible to compare benchmarks with an external tool :

```
go get -u golang.org/x/tools/cmd/benchcmp

go test ./fibonacci -bench . -benchmem > old.txt
(do some changes in the code)
go test ./fibonacci -bench . -benchmem > new.txt

~/go/bin/benchcmp old.txt new.txt
```

## Profiling

Profiling data are sampled and aggregated ones, not detailed traces. CPU profile measures elapsed time, and memory profile measures heap allocations (the stack is ignored).

While CPU benchmarks show how long an operation take (global view), profiling show which part of code consume time (detailed view). You get the same global/detailed view principle with memory.

### Profiling benchmarks

Get profiling data from the benchmarks:

* CPU profiling using `-cpuprofile=cpu.out`
* Memory profiling using `-benchmem -memprofile=mem.out`

An example with both :

```
go test ./fibonacci \
  -bench BenchmarkSuite \
  -benchmem \
  -cpuprofile=cpu.out \
  -memprofile=mem.out
```

CPU and memory profiling data from benchmarks are always stored in two separate files and will be analysed separately.

### Viewing profiling data

There are two way to exploit profiling data with the standard go tooling :

* through command line : `go tool pprof cpu.out`
* with a browser : `go tool pprof -http=localhost:8080 cpu.out`

The View menu :

* Top : ordered list of function sorted by their consumption of time/memory.
* Graph : function call tree, with time/memory annotations.
* Flamegraph : self-explanatory
* others...

### Profiling program with code

Use `pprof.StartCPUProfile()`, `pprof.StopCPUProfile()` and `pprof.WriteHeapProfile()`. See the `pprof` package documentation for more information.

A simple example :

```
cd showfib
go build
./showfib 30 2>cpu.out
go tool pprof -http=localhost:8080 cpu.out
```

## Tracing the GC work

This part isn't the most useful, but GC traces could easily spot a too high GC pressure.

Using an environment variable with any program (or test) : `GODEBUG=gctrace=1`

Or using code :

```go
  "runtime/trace"

	trace.Start(os.Stderr)
	defer trace.Stop()
```

Example with `webfib` :

```
cd webfib
go build
GODEBUG=gctrace=1 ./webfib

(other terminal)

go get -u github.com/rakyll/hey
~/go/bin/hey -n 1000 http://localhost:8000/?n=30
```

Graphical view with `gcvic`, which set `GODEBUG=gctrace=1` automatically :

```
go get -u github.com/davecheney/gcvis
~/go/bin/gcvis ./webfib

(other terminal)

~/go/bin/hey -n 2000 -c 200 http://localhost:8000/?n=30
```

In the previous example change the value of the `-c` argument and the `n` value to change the GC load and see the result.

## Tracing

Traces are event collected during the program execution. They give a chronological view of a program execution with detailed information about heap, GC, goroutines, core usage, ...

### With tests

Generate traces with a test, and visualize the data :

```
go test ./fibonacci \
  -bench BenchmarkSuite \
  -trace=trace.out

go tool trace trace.out
```

WARNING : Chrome is the only supported browser !

There are many detailed information... the blog post [Go execution tracer](https://blog.gopheracademy.com/advent-2017/go-execution-tracer/) is a good quick tour of the trace tool GUI.

Tip : from the *View trace* part hit `?` to show a help.


### With code in a program

Same principle as tracing with test, but you have to add code to collect traces into a file, and then use `go tool trace`.

```go
  "runtime/trace"

	trace.Start(os.Stderr)
	defer trace.Stop()
```

(no example, see `trace` package for more information)

## Tracing and profiling long running programs

Here Go tooling starts to really shine ! Go allows any program running a http server to be analysed during executing, even in production. And it's really easy.

Data are gathered only on demand, it's free when it is not used.

### Setup

Importing the `net/http/pprof` standard package add handler to `DefaultServeMux`. If there is already a web server using it, that's all.

```go
import _ "net/http/pprof"

// Add this only if needed
go func() {
	log.Println(http.ListenAndServe("localhost:8000", nil))
}()
```

The `webfib` exemple use it. Live data are available here : http://localhost:8000/debug/pprof/ , but is is it's not very user friendly.

### Security

Handlers provided by `net/http/pprof` should only be accessible to trusty client. It is not something you want to be available directly through internet, or internally by untrusted third party.

The package `net/http/pprof` register handlers to `DefaultServeMux`, use a separate server (create a dedicated `Mux`) for your http server. Each one using a different port, different security rules could be applied.  

### Profiling

It's possible to profile program using `net/http/pprof` with `go tool pprof`.

We'll use 3 terminals to :

* run the application
* collect traces (it could take one or two minutes)
* generate load

```
cd webfib
go build
./webfib

go tool pprof -http=localhost:8080 http://localhost:8000/

~/go/bin/hey -n 2000 -c 200 http://localhost:8000/?n=30
```

The `go tool ...` command line collect data, then open the browser. Be patient...

### Trace analysis

Tracing data have to be collected manually, and feed into `go tool trace`.

We'll use 3 terminals to :

* run application
* collect traces for 15 seconds
* generate load

```
cd webfib
go build
./webfib

curl -o trace.out http://localhost:8000/debug/pprof/trace?seconds=15

~/go/bin/hey -n 2000 -c 200 http://localhost:8000/?n=30
```

Now analyse the data with Chrome : `go tool trace trace.out`

## User-defined traces

User-defined traces were introduced with Go 1.11. It's not a new tool, but rather a way to add events to traces within your code.

What's in the box :

* `Task` struct : allows to trace high-level operations.
* `Region` struct : allow to trace lower-level operations.
* `Log` function : code can add log information to traces.

The `anowebfib` (like *another webfib*) use the three. The code identify each HTTP request with a `Task`, each call to fibonacci package with `Region`, and n are logged.

Example, with 3 terminals :

```
cd anowebfib
go build
./anowebfib

curl -o trace.out http://localhost:8000/debug/pprof/trace?seconds=20

~/go/bin/hey -n 2000 -c 200 http://localhost:8000/unique?n=30 
~/go/bin/hey -n 2000 -c 200 http://localhost:8000/multiple?n=30
```

As usual : `go tool trace trace.out`

See :
* User-defined tasks (and go down).
* User-defined regions (and go down).
* View trace : see events in goroutines.

## Useful links

Packages :

* Go testing package : https://golang.org/pkg/testing/
* Go runtime package : https://golang.org/pkg/runtime/
* Go trace package : https://golang.org/pkg/runtime/trace/
* Go pprof package : https://golang.org/pkg/runtime/pprof/

Others :

* How to Write Benchmarks in Go : https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go
* Profiling Go programs : https://blog.golang.org/profiling-go-programs
* Debugging performance issues in Go programs : https://github.com/golang/go/wiki/Performance
* Go execution tracer : https://blog.gopheracademy.com/advent-2017/go-execution-tracer/ (see also the *The tracer design doc* link)
* A whirlwind tour of Go’s runtime environment variables (see godebug) : https://dave.cheney.net/2015/11/29/a-whirlwind-tour-of-gos-runtime-environment-variables
* benchstat : https://godoc.org/golang.org/x/perf/cmd/benchstat
