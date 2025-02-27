package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"math/bits"
	"os"
	"os/user"
	"slices"
	"sort"
	"strconv"
	"syscall"
)

type FileType uint8

const (
	NonRegular FileType = iota
	RegularFile
	Directory
	Pipe
	SymbolicLink
	BlockDevice
	CharDevice
	Socket
)

type DisplayItem struct {
	Name        string
	Permissions string
	UserName    string
	GroupName   string
	ModifiedAt  string
	Size        int64
	NLinks      int
	Type        FileType
}

type DisplayItems []*DisplayItem

func (d DisplayItems) filterDirectories() DisplayItems {
	return slices.DeleteFunc(d, func(item *DisplayItem) bool {
		return item.Type != Directory
	})
}

func (d DisplayItems) filterFiles() DisplayItems {
	return slices.DeleteFunc(d, func(item *DisplayItem) bool {
		return item.Type == Directory
	})
}

type ByDirs []*DisplayItem

func (d ByDirs) Len() int {
	return len(d)
}

func (d ByDirs) Less(i, j int) bool {
	if d[i].Type == Directory && d[j].Type != Directory {
		return true
	}
	if d[i].Type != Directory && d[j].Type == Directory {
		return false
	}

	return d[i].Name < d[j].Name
}

func (d ByDirs) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

type ByFiles []*DisplayItem

func (d ByFiles) Len() int {
	return len(d)
}

func (d ByFiles) Less(i, j int) bool {
	if d[i].Type != Directory && d[j].Type == Directory {
		return true
	}
	if d[i].Type == Directory && d[j].Type != Directory {
		return false
	}

	return d[i].Name < d[j].Name
}

func (d ByFiles) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

type ColumnsWidth struct {
	LenPermissions int
	LenNLinks      int
	LenUserName    int
	LenGroupName   int
	LenSize        int
	LenDate        int
	LenFileName    int
}

func newColumnsWidth() *ColumnsWidth {
	return &ColumnsWidth{
		LenPermissions: 0,
		LenNLinks:      0,
		LenUserName:    0,
		LenGroupName:   0,
		LenSize:        0,
		LenDate:        0,
		LenFileName:    0,
	}
}

func (cw *ColumnsWidth) update(fileInfo fs.FileInfo, stat *syscall.Stat_t, userInfo *user.User, groupInfo *user.Group, config *Config) {
	if config.Filter.Permissions && len(fileInfo.Mode().String()) > cw.LenPermissions {
		cw.LenPermissions = len(fileInfo.Mode().String())
	}
	if config.Filter.NLinks && len(strconv.Itoa(int(stat.Nlink))) > cw.LenNLinks {
		cw.LenNLinks = len(strconv.Itoa(int(stat.Nlink)))
	}
	if config.Filter.UserName && len(userInfo.Username) > cw.LenUserName {
		cw.LenUserName = len(userInfo.Username)
	}
	if config.Filter.GroupName && len(groupInfo.Name) > cw.LenGroupName {
		cw.LenGroupName = len(groupInfo.Name)
	}
	if config.General.SizeUnit != None && len(SizeFormat(fileInfo.Size(), config.General.SizeUnit)) > cw.LenSize {
		cw.LenSize = len(SizeFormat(fileInfo.Size(), Auto))
	}
	if config.Filter.ModificationTime && len(fileInfo.ModTime().Format(config.General.DateFormat)) > cw.LenDate {
		cw.LenDate = len(fileInfo.ModTime().Format(config.General.DateFormat))
	}
	if config.Filter.FileName && len(fileInfo.Name()) > cw.LenFileName {
		cw.LenFileName = len(fileInfo.Name())
	}
}

func GetFiles(path string, config *Config) ([]*DisplayItem, *ColumnsWidth) {
	files, _ := os.ReadDir(path)

	listOfFiles := make([]*DisplayItem, 0, len(files))
	columnsWidth := newColumnsWidth()

	for _, f := range files {
		fileInfo, _ := f.Info()
		stat, _ := fileInfo.Sys().(*syscall.Stat_t)
		userInfo, _ := user.LookupId(strconv.Itoa(int(stat.Uid)))
		groupInfo, _ := user.LookupGroupId(strconv.Itoa(int(stat.Gid)))

		columnsWidth.update(fileInfo, stat, userInfo, groupInfo, config)
		fileType := typeOfFile(fileInfo)

		listOfFiles = append(listOfFiles, &DisplayItem{
			Name:        fileInfo.Name(),
			Permissions: fileInfo.Mode().String(),
			UserName:    userInfo.Username,
			GroupName:   groupInfo.Name,
			Size:        fileInfo.Size(),
			NLinks:      int(stat.Nlink),
			Type:        fileType,
			ModifiedAt:  fileInfo.ModTime().Format(config.General.DateFormat),
		})
	}

	if config.Filter.OnlyDirs {
		listOfFiles = DisplayItems(listOfFiles).filterDirectories()
	}
	if config.Filter.OnlyFiles {
		listOfFiles = DisplayItems(listOfFiles).filterFiles()
	}
	if config.General.DirsFirst {
		sort.Sort(ByDirs(listOfFiles))
	}
	if config.General.FilesFirst {
		sort.Sort(ByFiles(listOfFiles))
	}

	return listOfFiles, columnsWidth
}

func typeOfFile(fileInfo fs.FileInfo) FileType {
	if fileInfo.Mode().IsRegular() {
		return RegularFile
	}
	if fileInfo.Mode().Type() == fs.ModeDir {
		return Directory
	}
	if fileInfo.Mode().Type() == fs.ModeNamedPipe {
		return Pipe
	}
	if fileInfo.Mode().Type() == fs.ModeSymlink {
		return SymbolicLink
	}
	if fileInfo.Mode().Type() == fs.ModeDevice {
		return BlockDevice
	}
	if fileInfo.Mode().Type() == fs.ModeCharDevice {
		return CharDevice
	}
	if fileInfo.Mode().Type() == fs.ModeSocket {
		return Socket
	}

	return NonRegular
}

func SizeFormat(bytes int64, sizeType SizeType) string {
	if sizeType == Bytes {
		return fmt.Sprintf("%d B", bytes)
	}

	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	base := uint(bits.Len64(uint64(bytes)) / 10)
	val := float64(bytes) / float64(uint64(1<<(base*10)))

	return fmt.Sprintf("%.1f %ciB", val, " KMGTPE"[base])
}

func PathExists(path string) error {
	_, err := os.Stat(path)

	if err == nil {
		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return errors.New("path does not exist")
	}

	if errors.Is(err, os.ErrPermission) {
		return errors.New("permission denied")
	}

	return errors.New("unknown error: " + err.Error())
}
