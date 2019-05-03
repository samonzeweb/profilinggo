package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime/trace"
	"strconv"

	"github.com/samonzeweb/profilinggo/fibonacci"
)

func main() {
	http.HandleFunc("/unique", fibHandler)
	http.HandleFunc("/multiple", suiteHandler)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error on ListenAndServe : ", err)
	}
}

// Returns only the computed value for n
func fibHandler(w http.ResponseWriter, r *http.Request) {
	ctx, task := trace.NewTask(context.Background(), "UniqueFib")
	defer task.End()

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	n, ok := parseArgs(w, r)
	if !ok {
		return
	}

	trace.Log(ctx, "n value", strconv.Itoa(n))

	var result int
	trace.WithRegion(ctx, "Compute fibonacci (unique)", func() {
		result = fibonacci.Fibonacci(n)
	})
	fmt.Fprintf(w, "%d", result)
}

// Returns a table from 1 up to n
func suiteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, task := trace.NewTask(context.Background(), "MultipleFib")
	defer task.End()

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	n, ok := parseArgs(w, r)
	if !ok {
		return
	}

	trace.Log(ctx, "n value", strconv.Itoa(n))

	var result string
	trace.WithRegion(ctx, "Compute fibonacci suite", func() {
		result = fibonacci.Suite(n)
	})

	fmt.Fprintf(w, result)
}

// Parse the arguments
func parseArgs(w http.ResponseWriter, r *http.Request) (int, bool) {
	nStr := r.URL.Query().Get("n")
	n, err := strconv.Atoi(nStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing n query string parameter")
		return 0, false
	}
	if n < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "The value of n is invalid")
		return 0, false
	}

	return n, true
}
