package main

import (
	"fmt"

	"golang.org/x/tour/pic"
)

func Pic(dx, dy int) [][]uint8 {
	output := make([][]uint8, dy)

	for i := range output {
		output[i] = make([]uint8, dx)

		for xx := range output[i] {
			output[i][xx] = uint8((i ^ xx) / 2)
		}
	}

	return output
}

func main() {
	pic.Show(Pic)
	fmt.Printf("%v", Pic(255, 255))
}
