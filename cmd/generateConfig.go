package cmd

import (
	config2 "github.com/CezaryMackowski/go-ls/internal"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var configPath = ""

var generateConfigCmd = &cobra.Command{
	Use:          "generate",
	Short:        "Command used to generate default config",
	Example:      "go-ls generate",
	Version:      "1.0.0",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		pathBaseDir := path.Dir(configPath)

		err := os.MkdirAll(pathBaseDir, 0755)
		if err != nil {
			return err
		}

		configContent, err := toml.Marshal(config2.NewConfig())
		if err != nil {
			return err
		}

		err = os.WriteFile(configPath, configContent, 0755)
		if err != nil {
			return err
		}

		return err
	},
}

func init() {
	generateConfigCmd.Flags().StringVarP(&configPath, "path", "p", "~/.config/go-ls/config.toml", "")
}
