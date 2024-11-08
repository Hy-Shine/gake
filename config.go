package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var errNotSupportFileType = errors.New("don't support config file type")

type config struct {
	Platform  configPlatform `json:"platform" yaml:"platform"`
	OutputDir string         // ./output-20060102
	Targets   configTargets  `json:"targets" yaml:"targets"`
	// Env is env for compile
	//
	// eg: common: ["CGO_ENABLED=0"]
	Env configEnvArg `json:"env" yaml:"env"`

	// Args is args for compile
	//
	// eg: common: ["-trimpath"]
	Args        configEnvArg `json:"args" yaml:"args"`
	CompileCost bool         `json:"compileCost" yaml:"compileCost"`
	// FailSkip skip compile when compile error
	FailSkip bool `json:"failSkip" yaml:"failSkip"`
	// SuccessLog print success log
	SuccessLog bool `json:"successLog" yaml:"successLog"`
}

type configPlatform struct {
	OS      []string `json:"os" yaml:"os"`           // compile os
	Arch    []string `json:"arch" yaml:"arch"`       // compile arch
	Exclude []string `json:"exclude" yaml:"exclude"` // exclude os/arch
}

type configTarget struct {
	Entrance   string            `json:"entrance" yaml:"entrance"`
	OutputName string            `json:"outputName" yaml:"outputName"`
	Suffix     map[string]string `json:"suffix" yaml:"suffix"`
}

type configTargets struct {
	Suffix map[string]string `json:"suffix" yaml:"suffix"`
	Apps   []configTarget    `json:"apps" yaml:"apps"`
}

type configCompileArgs struct {
	BuildArgs []string `json:"buildArgs" yaml:"buildArgs"`
}

type configPlatformBase struct {
	// Use env/args for specific platform
	Use []string `json:"use" yaml:"use"`
	// Exclude env/args for specific platform, even if the arg in common list
	Exclude []string `json:"exclude" yaml:"exclude"`
}

type configEnvArg struct {
	Common   []string                      `json:"common" yaml:"common"`
	Platform map[string]configPlatformBase `json:"platform" yaml:"platform"`
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
	if len(cfg.Targets.Apps) == 0 {
		return errors.New("targets is empty")
	}
	if len(cfg.Platform.OS) == 0 {
		cfg.Platform.OS = []string{runtime.GOOS}
	}
	if len(cfg.Platform.Arch) == 0 {
		cfg.Platform.Arch = []string{runtime.GOARCH}
	}
	for i, target := range cfg.Targets.Apps {
		if target.Entrance == "" {
			return errors.New("compile entrance is empty")
		}
		if name := strings.TrimSpace(target.OutputName); name == "" {
			return fmt.Errorf("apps[%d] output name is empty", i)
		} else {
			cfg.Targets.Apps[i].OutputName = name
		}
		if target.Suffix == nil {
			target.Suffix = map[string]string{}
		}
	}
	cfg.OutputDir = fmt.Sprintf("./output-%s", time.Now().Format("20060102"))

	return nil
}
