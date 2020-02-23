package main

import (
	"fmt"
	"time"
)

func main() {
	defer showTime()()
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n) // slow
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}
func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}
func showTime() func() {
	now := time.Now()
	fmt.Printf("start time is %s", now)
	return func() {
		fmt.Printf("use time %s", time.Since(now))
	}
}
func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
