package internal

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type SizeType string

func (s *SizeType) String() string {
	return string(*s)
}

func (s *SizeType) Set(s2 string) error {
	*s = SizeType(s2)
	return nil
}

func (s *SizeType) Type() string {
	return "sizeType"
}

const (
	None      SizeType = "none"
	Bytes              = "B"
	KibiByte           = "KiB"
	MebiBytes          = "MiB"
	GibiBytes          = "GiB"
	Auto               = "auto"
)

type Color struct {
	Foreground string `toml:"foreground"`
	Background string `toml:"background"`
}

type Permissions struct {
	EmptyColor Color `toml:"empty_color"`

	OwnerReadColor    Color `toml:"owner_read_color"`
	OwnerWriteColor   Color `toml:"owner_write_color"`
	OwnerExecuteColor Color `toml:"owner_execute_color"`

	GroupReadColor    Color `toml:"group_read_color"`
	GroupWriteColor   Color `toml:"group_write_color"`
	GroupExecuteColor Color `toml:"group_execute_color"`

	OthersReadColor    Color `toml:"others_read_color"`
	OthersWriteColor   Color `toml:"others_write_color"`
	OthersExecuteColor Color `toml:"others_execute_color"`
}

type FileName struct {
	NonRegularColor   Color `toml:"non_regular_file_color"`
	RegularColor      Color `toml:"file_color"`
	DirectoryColor    Color `toml:"directory_color"`
	PipeColor         Color `toml:"pipe_color"`
	SymbolicLinkColor Color `toml:"symbolic_link_color"`
	BlockDeviceColor  Color `toml:"block_device_color"`
	CharDeviceColor   Color `toml:"char_device_color"`
	SocketColor       Color `toml:"socket_color"`
}

type UserName struct {
	ForegroundColor string `toml:"foreground_color"`
	BackgroundColor string `toml:"background_color"`
}

type General struct {
	Long       bool     `toml:"long"`
	DirsFirst  bool     `toml:"dirs_first"`
	FilesFirst bool     `toml:"files_first"`
	DateFormat string   `toml:"date_format"`
	SizeUnit   SizeType `toml:"size_unit"`
}

type Filter struct {
	FileName         bool `toml:"file_name"`
	Permissions      bool `toml:"permissions"`
	UserName         bool `toml:"user_name"`
	GroupName        bool `toml:"group_name"`
	ModificationTime bool `toml:"modification_time"`
	NLinks           bool `toml:"n_links"`
	All              bool `toml:"all"`
	OnlyDirs         bool `toml:"only_dirs"`
	OnlyFiles        bool `toml:"only_files"`
}

type Theme struct {
	NLinks           Color        `toml:"n_links"`
	UserName         Color        `toml:"user_name"`
	GroupName        Color        `toml:"group_name"`
	Size             Color        `toml:"size"`
	ModificationTime Color        `toml:"modification_time"`
	Permissions      *Permissions `toml:"permissions"`
	FileName         *FileName    `toml:"file_name"`
}

type Config struct {
	General *General `toml:"general"`
	Filter  *Filter  `toml:"filter"`
	Theme   *Theme   `toml:"theme"`
}

func NewConfig() *Config {
	return &Config{
		General: &General{
			Long:       true,
			DirsFirst:  true,
			FilesFirst: false,
			DateFormat: "Jan 02 15:04",
			SizeUnit:   Auto,
		},
		Filter: &Filter{
			FileName:         true,
			Permissions:      true,
			UserName:         true,
			GroupName:        true,
			ModificationTime: true,
			NLinks:           true,
			All:              true,
			OnlyDirs:         false,
			OnlyFiles:        false,
		},
		Theme: &Theme{
			NLinks:           Color{Foreground: "#D9D9D9", Background: ""},
			UserName:         Color{Foreground: "#EAEAC6", Background: ""},
			GroupName:        Color{Foreground: "#D4D584", Background: ""},
			Size:             Color{Foreground: "#FAF9D3", Background: ""},
			ModificationTime: Color{Foreground: "#00FF02", Background: ""},
			Permissions: &Permissions{
				EmptyColor:         Color{Foreground: "#C67D7D", Background: "#1825FF"},
				OwnerReadColor:     Color{Foreground: "#55BE57", Background: ""},
				OwnerWriteColor:    Color{Foreground: "#C1C27B", Background: ""},
				OwnerExecuteColor:  Color{Foreground: "#F4005F", Background: ""},
				GroupReadColor:     Color{Foreground: "#55BE57", Background: ""},
				GroupWriteColor:    Color{Foreground: "#C1C27B", Background: ""},
				GroupExecuteColor:  Color{Foreground: "#F4005F", Background: ""},
				OthersReadColor:    Color{Foreground: "#55BE57", Background: ""},
				OthersWriteColor:   Color{Foreground: "#C1C27B", Background: ""},
				OthersExecuteColor: Color{Foreground: "#F4005F", Background: ""},
			},
			FileName: &FileName{
				NonRegularColor:   Color{Foreground: "#FC971E", Background: ""},
				RegularColor:      Color{Foreground: "#FC971E", Background: ""},
				DirectoryColor:    Color{Foreground: "#05AEFF", Background: ""},
				PipeColor:         Color{Foreground: "#FC971E", Background: ""},
				SymbolicLinkColor: Color{Foreground: "#D71CFF", Background: ""},
				BlockDeviceColor:  Color{Foreground: "#FC971E", Background: ""},
				CharDeviceColor:   Color{Foreground: "#FC971E", Background: ""},
				SocketColor:       Color{Foreground: "#FC971E", Background: ""},
			},
		},
	}
}

func ParseConfigFile(path string) (*Config, error) {
	_, err := PathExists(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config *Config
	if err = toml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
