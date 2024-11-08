package main

import (
	"flag"
	"fmt"
)

func main() {
	var configPath string
	var help bool
	flag.BoolVar(&help, "h", false, "print help")
	flag.StringVar(&configPath, "c", "config.json", "config path")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	// read and check config
	cfg, err := initConfig(configPath)
	if err != nil {
		fmt.Printf("read config error: %v\n", err)
		return
	}

	err = compileBy(cfg)
	if err != nil {
		fmt.Printf("compile error: %v", err)
	}
}
