package main

import (
	"flag"
	"log"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.json", "config path")
	flag.Parse()

	err := compileBy(configPath)
	if err != nil {
		log.Printf("compile error: %v", err)
	}
}
