package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"womm/internal/checker"
	"womm/internal/config"

	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "womm",
	Short: "Works On My Machine - Environment checker and installer",
	Run: func(cmd *cobra.Command, args []string) {
		restore, _ := termenv.EnableVirtualTerminalProcessing(termenv.DefaultOutput())
		if restore != nil {
			defer restore()
		}

		resolvedConfigFile, ok := resolveConfigFile(cmd)
		if !ok {
			return
		}

		cfg, err := config.Load(resolvedConfigFile)
		if err != nil {
			fmt.Println("Error parsing", resolvedConfigFile, ":", err)
			return
		}

		runAll(cfg)
	},
}

func runAll(cfg *config.Config) {
	runChecks(cfg)

	tools := cfg.Dependencies
	if len(tools) == 0 {
		tools = cfg.Tools
	}

	if len(tools) > 0 {
		missing := checker.GetMissing(checker.CheckDeps(tools))
		if len(missing) > 0 {
			checker.PromptInstall(missing)
		}
	}

	if len(cfg.Commands) > 0 {
		checker.RunCommands(cfg.Commands)
	}
}

func runChecks(cfg *config.Config) {
	tools := cfg.Dependencies
	if len(tools) == 0 {
		tools = cfg.Tools
	}

	if len(cfg.Platforms) > 0 {
		supported := checker.CheckPlatforms(cfg.Platforms)
		checker.PrintPlatformResult(supported)
		if !supported {
			return
		}
	}

	if len(tools) > 0 {
		checker.PrintDepResults(checker.CheckDeps(tools))
	}

	if len(cfg.Env) > 0 {
		checker.PrintEnvResults(checker.CheckEnv(cfg.Env))
	}

	if len(cfg.Files) > 0 {
		checker.PrintFileResults(checker.CheckFiles(cfg.Files))
	}
}

func resolveConfigFile(cmd *cobra.Command) (string, bool) {
	if cmd.Flags().Changed("config") {
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			fmt.Printf("Config file %q was not found.\n", configFile)
			return "", false
		}
		return configFile, true
	}

	if _, err := os.Stat(".womm"); err == nil {
		return ".womm", true
	}

	matches, err := filepath.Glob("*.womm")
	if err != nil {
		fmt.Println("Error searching for .womm files:", err)
		return "", false
	}

	if len(matches) == 1 {
		return matches[0], true
	}

	if len(matches) > 1 {
		fmt.Println("Multiple .womm files found. Please choose one with -c:")
		for _, match := range matches {
			fmt.Printf("  go run main.go -c %s\n", match)
		}
		return "", false
	}

	fmt.Println("No .womm file found. Create .womm or provide a file path using -c.")
	return "", false
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", ".womm", "config file path")
	_ = rootCmd.Execute()
}
