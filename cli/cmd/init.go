package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var exampleWomm = `dependencies:
  node:
    version: ">=18.0"
    brew: node@18
    apt: nodejs
    choco: nodejs-lts

  git:
    version: ">=2.30"

  docker:
    version: ">=20.10"
    brew: --cask docker
    choco: docker-desktop
    winget: Docker.DockerDesktop

env:
  API_URL: http://localhost:3000
  MONGO_URI: mongodb://localhost:27017/app

files:
  - .env
  - config.json

commands:
  setup:
    - npm install
    - npm run build

platforms:
  windows: true
  linux: true
  macos: true
`

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [directory]",
	Short: "Create a .womm config file",
	Long: `Creates a .womm file with a detailed example configuration.

If a directory argument is given, the directory is created and .womm
is placed inside it. Otherwise .womm is written to the current directory.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := "."
		if len(args) > 0 {
			target = args[0]
			if err := os.MkdirAll(target, 0755); err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
		}

		path := filepath.Join(target, ".womm")
		if _, err := os.Stat(path); err == nil {
			fmt.Printf(".womm already exists at %s\n", path)
			return
		}

		if err := os.WriteFile(path, []byte(exampleWomm), 0644); err != nil {
			fmt.Println("Error writing .womm:", err)
			return
		}

		fmt.Println("Created", path)
	},
}
