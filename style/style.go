package style

import (
	"github.com/CezaryMackowski/go-ls/internal"
	"github.com/charmbracelet/lipgloss"
)

func PrintShortOutput(file *internal.DisplayItem, config *internal.Config) string {
	fgColor, bgColor := getFileTypeColor(file.Type, config)

	fileName := lipgloss.NewStyle().
		Width(len(file.Name)).
		Align(lipgloss.Left).
		MarginLeft(3).
		Foreground(fgColor).
		Background(bgColor).
		Render(file.Name)

	return fileName
}

func PrintLongOutput(file *internal.DisplayItem, config *internal.Config, columnsWidth *internal.ColumnsWidth) string {
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

	return lipgloss.JoinHorizontal(lipgloss.Left, permissions, nLinks, user, group, size, modifiedAt, fileName)
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
