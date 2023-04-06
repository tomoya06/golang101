package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	a, b := 0, 0

	return func() int {
		if b == 0 {
			b = 1
		} else {
			nb := a + b
			a = b
			b = nb
		}

		return a
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
