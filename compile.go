package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

func compileBy(c string) error {
	cfg, err := initConfig(c)
	if err != nil {
		return fmt.Errorf("read config error: %v", err)
	}

	// check output dir
	err = os.MkdirAll(cfg.Target.OutputDir, 0o755)
	if err != nil {
		return fmt.Errorf("create output dir error: %v", err)
	}

	// compile
	osList := removeDuplication(cfg.Platform.OS)
	archList := removeDuplication(cfg.Platform.Arch)
	for _, os := range osList {
		for _, arch := range archList {
			if contains(cfg.Platform.ExcludePlatform, combine(os, arch)) {
				continue
			}
			compileCfg := compileConfig{
				os:       os,
				arch:     arch,
				args:     cfg.CompileArgs.BuildArgs,
				entrance: cfg.Target.Entrance,
			}
			output := outputName(cfg.Target.OutputName, os, arch, cfg.Target.Suffix)
			compileCfg.output = path.Join(cfg.Target.OutputDir, output)
			err = compileByCmd(compileCfg)
			if err != nil {
				if !cfg.FailSkip {
					return err
				}
				log.Printf("compile error: %v", err)
				continue
			}
			if cfg.SuccessLog {
				fmt.Printf("compile success: %s in dir %s\n", output, cfg.Target.OutputDir)
			}
		}
	}

	return nil
}

type compileConfig struct {
	os, arch string
	entrance string
	output   string
	args     []string
	env      []string
}

func compileByCmd(cfg compileConfig) error {
	var args []string
	args = append(args, "build", "-o", cfg.output)
	if len(cfg.args) > 0 {
		args = append(args, cfg.args...)
	}
	args = append(args, cfg.entrance)

	cmd := exec.Command("go", args...)
	if len(cfg.env) > 0 {
		cmd.Env = append(os.Environ(), cfg.env...)
	}

	bs, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s, %v", string(bs), err)
	}

	return nil
}

func combine(os, arch string) string {
	return os + "/" + arch
}

func outputName(prefix, os, arch string, suffix map[string]string) string {
	name := os + "_" + arch
	if r, ok := suffix[combine(os, arch)]; ok && r != "" {
		name = r
	}
	output := prefix + "_" + name
	if os == "windows" {
		output += ".exe"
	}
	return output
}
