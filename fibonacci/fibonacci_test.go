package fibonacci

import (
	"fmt"
	"testing"
)

var knownFibonacci = map[int]int{
	1: 1,
	2: 1,
	3: 2,
	4: 3,
	5: 5,
	6: 8,
}

func TestFibonacci(t *testing.T) {
	for n, result := range knownFibonacci {
		description := fmt.Sprintf("Fib(%d)", n)
		t.Run(description, func(t *testing.T) {
			fib := Fibonacci(n)
			if fib != result {
				t.Error(fmt.Sprintf("Expected Fibonnaci(%d) == %d, but was %d", n, result, fib))
			}
		})
	}
}

func TestSuite(t *testing.T) {
	expected := "Fib(1)\t= 1\nFib(2)\t= 1\nFib(3)\t= 2"
	suite := Suite(3)
	if suite != expected {
		t.Errorf("Expected [%s] but was [%s]", expected, suite)
	}
}

func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Fibonacci(20)
	}
}

func BenchmarkSuite(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Suite(20)
	}
}
