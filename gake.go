package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

type cmdArgs struct {
	Config     string `short:"c" long:"config" description:"config path"`
	FailSkip   bool   `long:"fail-skip" description:"ignore compile fail"`
	SuccessLog bool   `long:"success-log" description:"print success log"`
	Cost       bool   `long:"cost" description:"print compile cost, must be used with flag success-log"`
	Help       bool   `short:"h" long:"help" description:"Show this help message"`
}

func main() {
	var arg cmdArgs
	parser := flags.NewParser(&arg, flags.PrintErrors)
	_, err := parser.Parse()
	if err != nil {
		fmt.Printf("parse args error: %v\n", err)
		return
	}

	if arg.Help {
		parser.WriteHelp(os.Stdout)
		return
	}

	// read and check config
	cfg, err := initConfig(arg.Config)
	if err != nil {
		fmt.Printf("read config error: %v\n", err)
		return
	}

	cfg.CompileCost = arg.Cost
	cfg.FailSkip = arg.FailSkip
	cfg.SuccessLog = arg.SuccessLog
	err = compileBy(cfg)
	if err != nil {
		fmt.Printf("compile error: %v", err)
	}
}
