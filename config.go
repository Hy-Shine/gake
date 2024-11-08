package main

import (
	"encoding/json"
	"errors"
	"os"
	"runtime"
)

type config struct {
	Platform    configPlatform    `json:"platform"`
	Target      configTarget      `json:"target"`
	CompileArgs configCompileArgs `json:"compileArgs"`
	FailSkip    bool              `json:"failSkip"`
	SuccessLog  bool              `json:"successLog"`
}

type configPlatform struct {
	OS      []string `json:"os"`
	Arch    []string `json:"arch"`
	Exclude []string `json:"exclude"`
}

type configTarget struct {
	Entrance   string            `json:"entrance"`
	OutputDir  string            `json:"outputDir"`
	OutputName string            `json:"outputName"`
	Suffix     map[string]string `json:"suffix"`
}

type configCompileArgs struct {
	BuildArgs []string `json:"buildArgs"`
}

func initConfig(c string) (config, error) {
	cfg, err := readConfig(c)
	if err != nil {
		return config{}, err
	}

	err = configHandle(&cfg)
	return cfg, err
}

func readConfig(c string) (config, error) {
	bs, err := os.ReadFile(c)
	if err != nil {
		return config{}, err
	}

	var cfg config
	err = json.Unmarshal(bs, &cfg)
	return cfg, err
}

func configHandle(cfg *config) error {
	if len(cfg.Platform.OS) == 0 {
		cfg.Platform.OS = []string{runtime.GOOS}
	}
	if len(cfg.Platform.Arch) == 0 {
		cfg.Platform.Arch = []string{runtime.GOARCH}
	}
	if cfg.Target.OutputName == "" {
		cfg.Target.OutputName = "output"
	}
	if cfg.Target.OutputDir == "" {
		cfg.Target.OutputDir = "./output"
	}
	if cfg.Target.Entrance == "" {
		return errors.New("compile entrance is empty")
	}

	return nil
}
