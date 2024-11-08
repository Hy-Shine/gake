package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"
)

func compileBy(cfg config) error {
	// check output dir
	err := os.MkdirAll(cfg.OutputDir, 0o755)
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

			for _, target := range cfg.Targets.Apps {
				now := time.Now()
				envs := getEnvArgs(cfg.Env.Common, cfg.Env.Platform[osArch(os, arch)])
				name := outputName(target.OutputName, os, arch, cfg.Targets.Suffix, target.Suffix)
				compileCfg := compileConfig{
					cost:     cfg.CompileCost,
					args:     getEnvArgs(cfg.Args.Common, cfg.Args.Platform[osArch(os, arch)]),
					env:      getEnvs(os, arch, envs),
					entrance: target.Entrance,
					output:   path.Join(cfg.OutputDir, name),
				}

				err = compileByCmd(compileCfg)
				if err != nil {
					if !cfg.FailSkip {
						return err
					}
					log.Printf("compile error for %s: %v", name, err)
					continue
				}
				if cfg.SuccessLog {
					var cost string
					if cfg.CompileCost {
						cost = fmt.Sprintf(", cost: %.1fs", time.Since(now).Seconds())
					}
					log.Printf("compile success: %s in dir %s%s\n", name, cfg.OutputDir, cost)
				}
			}
		}
	}

	return nil
}

type compileConfig struct {
	cost     bool
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

func outputName(prefix, os, arch string, commonSuffix, platformSuffix map[string]string) string {
	name := os + "_" + arch
	if r, ok := platformSuffix[osArch(os, arch)]; ok && r != "" {
		name = r
	} else if r, ok := commonSuffix[osArch(os, arch)]; ok && r != "" {
		name = r
	}
	output := prefix + "_" + name
	if os == "windows" {
		output += ".exe"
	}
	return output
}

func getEnvArgs(common []string, pf configPlatformBase) []string {
	common = append(common, pf.Use...)
	var commonEnv []string
	for _, e := range common {
		if !contains(pf.Exclude, e) {
			commonEnv = append(commonEnv, e)
		}
	}
	return removeDuplication(commonEnv)
}

func getEnvs(os, arch string, envs []string) []string {
	envs = append(envs, "GOOS="+os, "GOARCH="+arch)
	return removeDuplication(envs)
}
