package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	var count int
	flag.IntVar(&count, "c", 5, "count")
	flag.Parse()

	if count <= 0 || count > 60 {
		count = 10
	}

	for i := 1; i <= count; i++ {
		fmt.Printf("output: %d\n", i)
		time.Sleep(500 * time.Millisecond)
	}
}
