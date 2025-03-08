package cmd

import (
	"fmt"
	"github.com/CezaryMackowski/go-ls/internal"
	"github.com/CezaryMackowski/go-ls/style"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"os"
)

var (
	config           = internal.NewConfig()
	configFile       string
	isLong           bool
	dirsFirst        bool
	filesFirst       bool
	fileName         bool
	permissions      bool
	userName         bool
	groupName        bool
	modificationTime bool
	nLinks           bool
	all              bool
	onlyDirs         bool
	onlyFiles        bool
	dateFormat       string
	sizeUnit         internal.SizeType
)

var RootCmd = &cobra.Command{
	Use:          "go-ls ",
	Short:        "go-ls is a CLI tool to list files in local environment",
	Example:      "go-ls",
	Version:      "1.0.0",
	SilenceUsage: true,
	Args:         argsParse,
	PreRunE:      preRun,
	RunE:         run,
}

func init() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(generateConfigCmd)

	RootCmd.Flags().StringVarP(&configFile, "config-file", "c", "~/.config/go-ls/config.toml", "")
	RootCmd.Flags().StringVarP(&dateFormat, "date-format", "t", "Jan 02 15:04", "")
	RootCmd.Flags().VarP(&sizeUnit, "size-unit", "s", "")
	RootCmd.Flags().BoolVarP(&isLong, "long", "l", false, "")
	RootCmd.Flags().BoolVarP(&dirsFirst, "dirs-first", "d", true, "")
	RootCmd.Flags().BoolVarP(&filesFirst, "files-first", "f", false, "")
	RootCmd.Flags().BoolVarP(&fileName, "filename", "n", true, "")
	RootCmd.Flags().BoolVarP(&permissions, "permissions", "p", true, "")
	RootCmd.Flags().BoolVarP(&userName, "username", "u", true, "")
	RootCmd.Flags().BoolVarP(&groupName, "groupname", "g", true, "")
	RootCmd.Flags().BoolVarP(&modificationTime, "modification-time", "m", true, "")
	RootCmd.Flags().BoolVarP(&nLinks, "n-links", "r", true, "")
	RootCmd.Flags().BoolVarP(&all, "all", "a", true, "")
	RootCmd.Flags().BoolVarP(&onlyDirs, "only-dirs", "", false, "")
	RootCmd.Flags().BoolVarP(&onlyFiles, "only-files", "", false, "")
	RootCmd.MarkFlagsMutuallyExclusive("dirs-first", "files-first")
	RootCmd.MarkFlagsMutuallyExclusive("only-dirs", "only-files")
	RootCmd.SetErrPrefix("go-ls:")
}

func argsParse(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}

	if _, err := internal.PathExists(args[0]); err != nil {
		return err
	}

	return nil
}

func preRun(cmd *cobra.Command, _ []string) error {
	var err error
	if cmd.Flags().Changed("config-file") {
		config, err = internal.ParseConfigFile(configFile)
		if err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("date-format") {
		config.General.DateFormat = dateFormat
	}
	if cmd.Flags().Changed("size-unit") {
		config.General.SizeUnit = sizeUnit
	}
	if cmd.Flags().Changed("long") {
		config.General.Long = isLong
	}
	if cmd.Flags().Changed("dirs-first") {
		config.General.DirsFirst = dirsFirst
	}
	if cmd.Flags().Changed("files-first") {
		config.General.FilesFirst = filesFirst
	}
	if cmd.Flags().Changed("filename") {
		config.Filter.FileName = fileName
	}
	if cmd.Flags().Changed("permissions") {
		config.Filter.Permissions = permissions
	}
	if cmd.Flags().Changed("username") {
		config.Filter.UserName = userName
	}
	if cmd.Flags().Changed("groupname") {
		config.Filter.GroupName = groupName
	}
	if cmd.Flags().Changed("modification-time") {
		config.Filter.ModificationTime = modificationTime
	}
	if cmd.Flags().Changed("n-links") {
		config.Filter.NLinks = nLinks
	}
	if cmd.Flags().Changed("all") {
		config.Filter.All = all
	}
	if cmd.Flags().Changed("only-dirs") {
		config.Filter.OnlyDirs = onlyDirs
	}
	if cmd.Flags().Changed("only-files") {
		config.Filter.OnlyFiles = onlyFiles
	}

	return nil
}

func run(_ *cobra.Command, args []string) error {
	var lines []string
	var output string
	files, columnsWidth, err := internal.GetFiles(args[0], config)
	if err != nil {
		return err
	}

	for _, f := range files {
		if config.General.Long {
			output = style.PrintLongOutput(f, config, columnsWidth)
		} else {
			output = style.PrintShortOutput(f, config)
		}

		lines = append(lines, output)
	}

	if config.General.Long {
		fmt.Println(lipgloss.JoinVertical(lipgloss.Top, lines...))
	} else {
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Left, lines...))
	}

	return nil
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
