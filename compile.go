package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

func compileBy(cfg config) error {
	// check output dir
	err := os.MkdirAll(cfg.Target.OutputDir, 0o755)
	if err != nil {
		return fmt.Errorf("create output dir error: %v", err)
	}

	// compile
	osList := removeDuplication(cfg.Platform.OS)
	archList := removeDuplication(cfg.Platform.Arch)
	exclude := removeDuplication(cfg.Platform.Exclude)
	for _, os := range osList {
		for _, arch := range archList {
			if contains(exclude, osArch(os, arch)) {
				continue
			}

			compileCfg := compileConfig{
				args:     cfg.CompileArgs.BuildArgs,
				env:      []string{"GOOS=" + os, "GOARCH=" + arch},
				entrance: cfg.Target.Entrance,
			}
			output := outputName(cfg.Target.OutputName, os, arch, cfg.Target.Suffix)
			compileCfg.output = path.Join(cfg.Target.OutputDir, output)
			err = compileByCmd(compileCfg)
			if err != nil {
				if !cfg.FailSkip {
					return err
				}
				log.Printf("compile error for %s: %v", output, err)
				continue
			}
			if cfg.SuccessLog {
				log.Printf("compile success: %s in dir %s\n", output, cfg.Target.OutputDir)
			}
		}
	}

	return nil
}

type compileConfig struct {
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
		if len(bs) == 0 {
			return err
		}
		return fmt.Errorf("output: %s, err: %v", string(bs), err)
	}

	return nil
}

func osArch(os, arch string) string {
	return os + "/" + arch
}

func outputName(prefix, os, arch string, suffix map[string]string) string {
	name := os + "_" + arch
	if r, ok := suffix[osArch(os, arch)]; ok && r != "" {
		name = r
	}
	output := prefix + "_" + name
	if os == "windows" {
		output += ".exe"
	}
	return output
}
