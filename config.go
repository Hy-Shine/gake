package main

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"runtime"

	"gopkg.in/yaml.v2"
)

var errNotSupportFileType = errors.New("don't support config file type")

type config struct {
	Platform    configPlatform    `json:"platform" yaml:"platform"`
	Target      configTarget      `json:"target" yaml:"target"`
	CompileArgs configCompileArgs `json:"compileArgs" yaml:"compileArgs"`
	FailSkip    bool              `json:"failSkip" yaml:"failSkip"`
	SuccessLog  bool              `json:"successLog" yaml:"successLog"`
}

type configPlatform struct {
	OS      []string `json:"os" yaml:"os"`
	Arch    []string `json:"arch" yaml:"arch"`
	Exclude []string `json:"exclude" yaml:"exclude"`
}

type configTarget struct {
	Entrance   string            `json:"entrance" yaml:"entrance"`
	OutputDir  string            `json:"outputDir" yaml:"outputDir"`
	OutputName string            `json:"outputName" yaml:"outputName"`
	Suffix     map[string]string `json:"suffix" yaml:"suffix"`
}

type configCompileArgs struct {
	BuildArgs []string `json:"buildArgs" yaml:"buildArgs"`
}

func initConfig(c string) (config, error) {
	cfg, err := readConfig(c)
	if err != nil {
		return config{}, err
	}

	err = configHandle(&cfg)
	return cfg, err
}

func checkFileType(c string) (string, error) {
	ext := path.Ext(c)
	if ext == ".json" || ext == ".yaml" {
		return ext, nil
	}

	return "", errNotSupportFileType
}

func readConfig(c string) (config, error) {
	bs, err := os.ReadFile(c)
	if err != nil {
		return config{}, err
	}

	var cfg config
	_ft, err := checkFileType(c)
	if err != nil {
		return config{}, err
	}

	if _ft == ".json" {
		err = json.Unmarshal(bs, &cfg)
	} else {
		err = yaml.Unmarshal(bs, &cfg)
	}

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
