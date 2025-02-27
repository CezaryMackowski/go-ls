package cmd

import (
	"fmt"
	"github.com/CezaryMackowski/go-ls/internal"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"os"
	"strconv"
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
	cobra.OnInitialize(initConfig)

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
	RootCmd.SetFlagErrorFunc(errorFunc)
	RootCmd.SetErrPrefix("go-ls:")
}

func errorFunc(cmd *cobra.Command, err error) error {
	g := 2 + 2

	fmt.Println(g)

	return nil
}

func initConfig() {

}

func argsParse(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return err
	}

	if err := internal.PathExists(args[0]); err != nil {
		return err
	}

	return nil
}

func preRun(cmd *cobra.Command, args []string) error {
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

func run(cmd *cobra.Command, args []string) error {
	var lines []string
	var output string
	files, columnsWidth := internal.GetFiles(args[0], config)

	for _, f := range files {
		if config.General.Long {
			output = printLongOutput(f, config, columnsWidth)
		} else {
			output = printShortOutput(f, config)
		}

		lines = append(lines, output)
	}

	// TODO: Prawdopodobnie będzie trzeba zmienić po implementacji printShortOutput
	fmt.Println(lipgloss.JoinVertical(lipgloss.Top, lines...))

	return nil
}

func printShortOutput(file *internal.DisplayItem, config *internal.Config) string {
	// TODO: Zaimplementować całą funkcje

	return ""
}

func printLongOutput(file *internal.DisplayItem, config *internal.Config, columnsWidth *internal.ColumnsWidth) string {
	var permissions, nLinks, user, group, size, modifiedAt, fileName string

	if columnsWidth.LenPermissions != 0 {
		permissions = formatPermissions(file, config, columnsWidth.LenPermissions)
	}
	if columnsWidth.LenNLinks != 0 {
		nLinks = formatNLinks(file, config, columnsWidth.LenNLinks)
	}
	if columnsWidth.LenUserName != 0 {
		user = formatUser(file, config, columnsWidth.LenUserName)
	}
	if columnsWidth.LenGroupName != 0 {
		group = formatGroup(file, config, columnsWidth.LenGroupName)
	}
	if columnsWidth.LenSize != 0 {
		size = formatSize(file, config, columnsWidth.LenSize)
	}
	if columnsWidth.LenDate != 0 {
		modifiedAt = formatModifiedAt(file, config, columnsWidth.LenDate)
	}
	if columnsWidth.LenFileName != 0 {
		fileName = formatFileName(file, config, columnsWidth.LenFileName)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, permissions, nLinks, user, group, size, modifiedAt, fileName)
}

func formatPermissions(file *internal.DisplayItem, config *internal.Config, columnWidth int) string {
	return lipgloss.NewStyle().
		Width(columnWidth + 2).
		Align(lipgloss.Right).
		Render(lipgloss.JoinHorizontal(
			lipgloss.Right,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.DirectoryColor.ForegroundColor)).
				Background(lipgloss.Color(config.Theme.Permissions.DirectoryColor.BackgroundColor)).
				Render(string(file.Permissions[0])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.OwnerReadColor.ForegroundColor)).
				Background(lipgloss.Color(config.Theme.Permissions.OwnerReadColor.BackgroundColor)).
				Render(string(file.Permissions[1])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.OwnerWriteColor.ForegroundColor)).
				Background(lipgloss.Color(config.Theme.Permissions.OwnerWriteColor.BackgroundColor)).
				Render(string(file.Permissions[2])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.OwnerExecuteColor.ForegroundColor)).
				Background(lipgloss.Color(config.Theme.Permissions.OwnerExecuteColor.BackgroundColor)).
				Render(string(file.Permissions[3])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.GroupReadColor.ForegroundColor)).
				Background(lipgloss.Color(config.Theme.Permissions.GroupReadColor.BackgroundColor)).
				Render(string(file.Permissions[4])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.GroupWriteColor.ForegroundColor)).
				Background(lipgloss.Color(config.Theme.Permissions.GroupWriteColor.BackgroundColor)).
				Render(string(file.Permissions[5])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.GroupExecuteColor.ForegroundColor)).
				Background(lipgloss.Color(config.Theme.Permissions.GroupExecuteColor.BackgroundColor)).
				Render(string(file.Permissions[6])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.OthersReadColor.ForegroundColor)).
				Background(lipgloss.Color(config.Theme.Permissions.OthersReadColor.BackgroundColor)).
				Render(string(file.Permissions[7])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.OthersWriteColor.ForegroundColor)).
				Background(lipgloss.Color(config.Theme.Permissions.OthersWriteColor.BackgroundColor)).
				Render(string(file.Permissions[8])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.OthersExecuteColor.ForegroundColor)).
				Background(lipgloss.Color(config.Theme.Permissions.OthersExecuteColor.BackgroundColor)).
				Render(string(file.Permissions[9])),
		),
		)
}

func formatFileName(file *internal.DisplayItem, config *internal.Config, columnWidth int) string {
	var fgColor, bgColor lipgloss.Color

	switch file.Type {
	case internal.RegularFile:
		fgColor = lipgloss.Color(config.Theme.FileName.RegularFileColor.ForegroundColor)
		bgColor = lipgloss.Color(config.Theme.FileName.RegularFileColor.BackgroundColor)
	case internal.Directory:
		fgColor = lipgloss.Color(config.Theme.FileName.DirectoryColor.ForegroundColor)
		bgColor = lipgloss.Color(config.Theme.FileName.DirectoryColor.BackgroundColor)
	case internal.Pipe:
		fgColor = lipgloss.Color(config.Theme.FileName.PipeColor.ForegroundColor)
		bgColor = lipgloss.Color(config.Theme.FileName.PipeColor.BackgroundColor)
	case internal.SymbolicLink:
		fgColor = lipgloss.Color(config.Theme.FileName.SymbolicLinkColor.ForegroundColor)
		bgColor = lipgloss.Color(config.Theme.FileName.SymbolicLinkColor.BackgroundColor)
	case internal.BlockDevice:
		fgColor = lipgloss.Color(config.Theme.FileName.BlockDeviceColor.ForegroundColor)
		bgColor = lipgloss.Color(config.Theme.FileName.BlockDeviceColor.BackgroundColor)
	case internal.CharDevice:
		fgColor = lipgloss.Color(config.Theme.FileName.CharDeviceColor.ForegroundColor)
		bgColor = lipgloss.Color(config.Theme.FileName.CharDeviceColor.BackgroundColor)
	case internal.Socket:
		fgColor = lipgloss.Color(config.Theme.FileName.SocketColor.ForegroundColor)
		bgColor = lipgloss.Color(config.Theme.FileName.SocketColor.BackgroundColor)
	default:
		fgColor = lipgloss.Color(config.Theme.FileName.NonRegularFileColor.ForegroundColor)
		bgColor = lipgloss.Color(config.Theme.FileName.NonRegularFileColor.BackgroundColor)
	}

	return lipgloss.NewStyle().
		Width(columnWidth + 2).
		Align(lipgloss.Left).
		MarginLeft(3).
		Foreground(fgColor).
		Background(bgColor).
		Render(file.Name)
}

func formatNLinks(file *internal.DisplayItem, config *internal.Config, columnWidth int) string {
	return lipgloss.NewStyle().
		Width(columnWidth + 2).
		Align(lipgloss.Right).
		MarginLeft(3).
		Foreground(lipgloss.Color(config.Theme.NLinks.ForegroundColor)).
		Background(lipgloss.Color(config.Theme.NLinks.BackgroundColor)).
		Render(strconv.Itoa(file.NLinks))
}

func formatUser(file *internal.DisplayItem, config *internal.Config, columnWidth int) string {
	return lipgloss.NewStyle().
		Width(columnWidth + 2).
		Align(lipgloss.Right).
		MarginLeft(3).
		Foreground(lipgloss.Color(config.Theme.UserName.ForegroundColor)).
		Background(lipgloss.Color(config.Theme.UserName.BackgroundColor)).
		Render(file.UserName)
}

func formatGroup(file *internal.DisplayItem, config *internal.Config, columnWidth int) string {
	return lipgloss.NewStyle().
		Width(columnWidth + 2).
		Align(lipgloss.Right).
		MarginLeft(3).
		Foreground(lipgloss.Color(config.Theme.GroupName.ForegroundColor)).
		Background(lipgloss.Color(config.Theme.GroupName.BackgroundColor)).
		Render(file.GroupName)
}

func formatSize(file *internal.DisplayItem, config *internal.Config, columnWidth int) string {
	return lipgloss.NewStyle().
		Width(columnWidth + 2).
		Align(lipgloss.Right).
		MarginLeft(3).
		Foreground(lipgloss.Color(config.Theme.UserName.ForegroundColor)).
		Background(lipgloss.Color(config.Theme.UserName.BackgroundColor)).
		Render(internal.SizeFormat(file.Size, config.General.SizeUnit))
}

func formatModifiedAt(file *internal.DisplayItem, config *internal.Config, columnWidth int) string {
	return lipgloss.NewStyle().
		Width(columnWidth + 2).
		Align(lipgloss.Left).
		MarginLeft(3).
		Foreground(lipgloss.Color(config.Theme.ModificationTime.ForegroundColor)).
		Background(lipgloss.Color(config.Theme.ModificationTime.BackgroundColor)).
		Render(file.ModifiedAt)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
