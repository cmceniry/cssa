package stash

import (
	"fmt"
	"github.com/cmceniry/cssa/util"
)

type Archive interface {
	AddFile(*util.SnapFile) error
	Commit() error
	GetFilelist() []string
	GetName() string
	IsValid() bool
	SetInvalid()
	SetValid()
}

type Stash interface {
	//AddFile(
	CreateArchive(string) Archive
	GetArchive(string) Archive
	GetArchivelist() []string
	GetFileSize(string) int64
	IsFileExist(*util.SnapFile) bool
	Save() error
}

func ParseConfig(c map[string]interface{}) (Stash, error) {
	val, ok := c["type"]
	if !ok {
		return nil, fmt.Errorf("No stash type identified")
	}
	v, ok := c["type"].(string)
	if !ok {
		return nil, fmt.Errorf("Invalid stash type value: %s(%T)", val, val)
	}
	switch v {
	case "hollow":
		s, err := NewHollowStashFromConfig(c)
		return s, err
	default:
		return nil, fmt.Errorf("Unknown stash type: %s", v)
	}
}

