package fibonacci

import "fmt"

// Fibonacci computes the n'th value of the Fibonacci suite
// (it starts with n=1)
func Fibonacci(n int) int {
	if n <= 2 {
		return 1
	}

	return Fibonacci(n-1) + Fibonacci(n-2)
}

// Suite computes the fibonacci suite upto
// the given argument and returns the result as
// a string.
// It's higly inefficient to give some work to the GC.
func Suite(n int) string {
	table := ""
	for i := 1; i <= n; i++ {
		f := Fibonacci(i)
		if len(table) > 0 {
			table = table + "\n"
		}
		table = table + fmt.Sprintf("Fib(%d)\t= %d", i, f)
	}
	return table
}
