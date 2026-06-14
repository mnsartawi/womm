package cmd

import (
	"fmt"
	"womm/internal/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(verifyCmd)
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Check environment against .womm without installing anything",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		resolved, ok := resolveConfigFile(cmd)
		if !ok {
			return
		}

		cfg, err := config.Load(resolved)
		if err != nil {
			fmt.Println("Error parsing", resolved, ":", err)
			return
		}

		runChecks(cfg)
	},
}
