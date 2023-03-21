package main

import (
	"fmt"
	"math"
)

func main() {
	count := 0.0
	for true {
		count += 1.0 * math.Pi
		fmt.Println("Hello, Shakti", count)
		//fmt.Println("Hello, Shakti")
		//time.Sleep(5 * time.Second)
	}
}
