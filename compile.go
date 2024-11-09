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
	osList := distinct(split(cfg.Platform.OS))
	archList := distinct(split(cfg.Platform.Arch))
	exclude := distinct(split(cfg.Platform.Exclude))
	for _, os := range osList {
		for _, arch := range archList {
			osArchStr := osArch(os, arch)
			if contains(exclude, osArchStr) {
				continue
			}

			for _, target := range cfg.Targets.Apps {
				envs := getEnvArgs(cfg.Env.Common, cfg.Env.Platform[os], cfg.Env.Platform[osArchStr])
				name := outputName(target.OutputName, os, arch, cfg.Targets.NameSuffix, target.NameSuffix)
				compileCfg := compileConfig{
					args:     getEnvArgs(cfg.Args.Common, cfg.Args.Platform[os], cfg.Args.Platform[osArchStr]),
					env:      getEnvs(os, arch, envs),
					entrance: target.Entrance,
					output:   path.Join(cfg.OutputDir, name),
				}

				now := time.Now()
				err = compileByCmd(compileCfg)
				if err != nil {
					if !cfg.FailSkip {
						return err
					}
					log.Printf("failed, name: %s, err: %v", name, err)
					continue
				}
				if cfg.SuccessLog {
					var cost string
					if cfg.CompileCost {
						cost = fmt.Sprintf(", cost: %.1fs", time.Since(now).Seconds())
					}
					log.Printf("success, dir: %s, name: %s%s\n", cfg.OutputDir, name, cost)
				}
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

func outputName(name, os, arch string, commonSuffix, platformSuffix map[string]string) string {
	osArchStr := osArch(os, arch)
	suffix := os + "_" + arch
	if r, ok := platformSuffix[osArchStr]; ok && r != "" {
		suffix = r
	} else if r, ok := commonSuffix[osArchStr]; ok && r != "" {
		suffix = r
	}
	output := name + "_" + suffix
	if os == "windows" {
		output += ".exe"
	}
	return output
}

func getEnvArgs(common []string, osConfig, pfConfig configPlatformBase) []string {
	var base configPlatformBase
	if len(pfConfig.Use) > 0 {
		base.Use = pfConfig.Use
	} else {
		base.Use = osConfig.Use
	}
	if len(pfConfig.Exclude) > 0 {
		base.Exclude = pfConfig.Exclude
	} else {
		base.Exclude = osConfig.Exclude
	}

	common = append(common, base.Use...)
	var commonEnv []string
	for _, e := range common {
		if !contains(base.Exclude, e) {
			commonEnv = append(commonEnv, e)
		}
	}
	return distinct(commonEnv)
}

func getEnvs(os, arch string, envs []string) []string {
	envs = append(envs, "GOOS="+os, "GOARCH="+arch)
	return distinct(envs)
}
