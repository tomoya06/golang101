package main

import (
	"fmt"
)

func Sqrt(x float64) float64 {
	z := 1.0

	for i := 0; i < 10; i++ {
		diff := (z*z - x) / (2 * z)
		fmt.Println(diff)
		z -= diff
	}

	return z
}

func main() {
	fmt.Println(Sqrt(2))
}
