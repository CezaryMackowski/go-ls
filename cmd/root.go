package cmd

import (
	"fmt"
	"github.com/CezaryMackowski/go-ls/internal"
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
		nLinks = formatCommonColumn(
			file.NLinks,
			columnsWidth.LenNLinks+2,
			lipgloss.Color(config.Theme.NLinks.Foreground),
			lipgloss.Color(config.Theme.NLinks.Background),
			lipgloss.Right,
		)
	}
	if columnsWidth.LenUserName != 0 {
		user = formatCommonColumn(
			file.UserName,
			columnsWidth.LenUserName+2,
			lipgloss.Color(config.Theme.NLinks.Foreground),
			lipgloss.Color(config.Theme.NLinks.Background),
			lipgloss.Right,
		)
	}
	if columnsWidth.LenGroupName != 0 {
		group = formatCommonColumn(
			file.GroupName,
			columnsWidth.LenGroupName+2,
			lipgloss.Color(config.Theme.GroupName.Foreground),
			lipgloss.Color(config.Theme.GroupName.Background),
			lipgloss.Right,
		)
	}
	if columnsWidth.LenSize != 0 {
		size = formatCommonColumn(
			internal.SizeFormat(file.Size, config.General.SizeUnit),
			columnsWidth.LenSize+2,
			lipgloss.Color(config.Theme.Size.Foreground),
			lipgloss.Color(config.Theme.Size.Background),
			lipgloss.Right,
		)
	}
	if columnsWidth.LenModifiedAt != 0 {
		modifiedAt = formatCommonColumn(
			file.ModifiedAt,
			columnsWidth.LenModifiedAt+2,
			lipgloss.Color(config.Theme.ModificationTime.Foreground),
			lipgloss.Color(config.Theme.ModificationTime.Background),
			lipgloss.Left,
		)
	}
	if columnsWidth.LenFileName != 0 {
		fgColor, bgColor := getFileTypeColor(file.Type, config)

		fileName = formatCommonColumn(
			file.Name,
			columnsWidth.LenFileName+2,
			fgColor,
			bgColor,
			lipgloss.Left,
		)
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
				Foreground(lipgloss.Color(config.Theme.Permissions.OwnerReadColor.Foreground)).
				Background(lipgloss.Color(config.Theme.Permissions.OwnerReadColor.Background)).
				Render(string(file.Permissions[1])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.OwnerWriteColor.Foreground)).
				Background(lipgloss.Color(config.Theme.Permissions.OwnerWriteColor.Background)).
				Render(string(file.Permissions[2])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.OwnerExecuteColor.Foreground)).
				Background(lipgloss.Color(config.Theme.Permissions.OwnerExecuteColor.Background)).
				Render(string(file.Permissions[3])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.GroupReadColor.Foreground)).
				Background(lipgloss.Color(config.Theme.Permissions.GroupReadColor.Background)).
				Render(string(file.Permissions[4])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.GroupWriteColor.Foreground)).
				Background(lipgloss.Color(config.Theme.Permissions.GroupWriteColor.Background)).
				Render(string(file.Permissions[5])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.GroupExecuteColor.Foreground)).
				Background(lipgloss.Color(config.Theme.Permissions.GroupExecuteColor.Background)).
				Render(string(file.Permissions[6])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.OthersReadColor.Foreground)).
				Background(lipgloss.Color(config.Theme.Permissions.OthersReadColor.Background)).
				Render(string(file.Permissions[7])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.OthersWriteColor.Foreground)).
				Background(lipgloss.Color(config.Theme.Permissions.OthersWriteColor.Background)).
				Render(string(file.Permissions[8])),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(config.Theme.Permissions.OthersExecuteColor.Foreground)).
				Background(lipgloss.Color(config.Theme.Permissions.OthersExecuteColor.Background)).
				Render(string(file.Permissions[9])),
		),
		)
}

func getNonExistingPermissionColor(permission string, config *internal.Config, defaultColor internal.Color) internal.Color {
	if permission == "-" {
		return internal.Color{
			Foreground: config.Theme.Permissions.EmptyColor.Foreground,
			Background: config.Theme.Permissions.EmptyColor.Background,
		}
	}

	return defaultColor
}

func formatCommonColumn(text string, width int, fgColor lipgloss.Color, bgColor lipgloss.Color, align lipgloss.Position) string {
	return lipgloss.NewStyle().
		Width(width).
		Align(align).
		MarginLeft(3).
		Foreground(fgColor).
		Background(bgColor).
		Render(text)
}

func getFileTypeColor(fileType internal.FileType, config *internal.Config) (lipgloss.Color, lipgloss.Color) {
	var fgColor, bgColor lipgloss.Color

	switch fileType {
	case internal.Regular:
		fgColor = lipgloss.Color(config.Theme.FileName.RegularColor.Foreground)
		bgColor = lipgloss.Color(config.Theme.FileName.RegularColor.Background)
	case internal.Directory:
		fgColor = lipgloss.Color(config.Theme.FileName.DirectoryColor.Foreground)
		bgColor = lipgloss.Color(config.Theme.FileName.DirectoryColor.Background)
	case internal.Pipe:
		fgColor = lipgloss.Color(config.Theme.FileName.PipeColor.Foreground)
		bgColor = lipgloss.Color(config.Theme.FileName.PipeColor.Background)
	case internal.SymbolicLink:
		fgColor = lipgloss.Color(config.Theme.FileName.SymbolicLinkColor.Foreground)
		bgColor = lipgloss.Color(config.Theme.FileName.SymbolicLinkColor.Background)
	case internal.BlockDevice:
		fgColor = lipgloss.Color(config.Theme.FileName.BlockDeviceColor.Foreground)
		bgColor = lipgloss.Color(config.Theme.FileName.BlockDeviceColor.Background)
	case internal.CharDevice:
		fgColor = lipgloss.Color(config.Theme.FileName.CharDeviceColor.Foreground)
		bgColor = lipgloss.Color(config.Theme.FileName.CharDeviceColor.Background)
	case internal.Socket:
		fgColor = lipgloss.Color(config.Theme.FileName.SocketColor.Foreground)
		bgColor = lipgloss.Color(config.Theme.FileName.SocketColor.Background)
	default:
		fgColor = lipgloss.Color(config.Theme.FileName.NonRegularColor.Foreground)
		bgColor = lipgloss.Color(config.Theme.FileName.NonRegularColor.Background)
	}

	return fgColor, bgColor
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
