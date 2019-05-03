package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strconv"

	"github.com/samonzeweb/profilinggo/fibonacci"
)

func main() {
	http.HandleFunc("/", fibHandler)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error on ListenAndServe : ", err)
	}
}

func fibHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	nStr := r.URL.Query().Get("n")
	n, err := strconv.Atoi(nStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing n query string parameter")
		return
	}
	if n < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "The value of n is invalid")
		return
	}
	fmt.Fprintf(w, fibonacci.Suite(n))
}
