package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:          "version",
	Short:        "Command used to get version information",
	Example:      "go-ls version",
	Version:      "1.0.0",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(cmd.Version)
		return nil
	},
}
