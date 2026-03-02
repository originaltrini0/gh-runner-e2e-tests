package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	fmt.Printf("Building and running Go binary on architecture: %s\n", runtime.GOARCH)

	// Simulate some computational work
	fmt.Println("Starting computational work...")
	rand.Seed(time.Now().UnixNano())

	iterations := 10000000
	sum := 0.0
	for i := 0; i < iterations; i++ {
		sum += rand.Float64()
	}

	fmt.Printf("Completed %d iterations. Result: %f\n", iterations, sum)
}
