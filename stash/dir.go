package stash

import (
	"github.com/cmceniry/cssa/util"
	"fmt"
	"os"
)

/*
 *
 * Directory backed Archive
 *
 */

type DirArchive struct {
}

func (a *DirArchive) AddFile(sf *util.SnapFile) error {
	// TODO
	return nil
}

func (a *DirArchive) Commit() error {
	// TODO
	return nil
}

func (a *DirArchive) GetFilelist() []string {
	// TODO
	return []string{}
}

func (a *DirArchive) GetName() string {
	// TODO
	return ""
}

func (a *DirArchive) IsValid() bool {
	// TODO
	return false
}

func (a *DirArchive) SetInvalid() {
	// TODO
}

func (a *DirArchive) SetValid() {
	// TODO
}

/*
 *
 * Directory backed Stash
 *
 */

type DirStash struct {
	InventoryDir	string
	Loaded		bool	`yaml:"-"`
}

func NewDirStashFromConfig(c map[string]interface{}) (*DirStash, error) {
	if dirnameintf, ok := c["inventorydir"]; ok {
		if dirname, ok := dirnameintf.(string); ok {
			fmt.Println(dirname)
			s := NewDirStash(dirname)
			return s, nil
		} else {
			return nil, fmt.Errorf("Invalid inventory directory: %s(%T)", dirnameintf, dirnameintf)
		}
	} else {
		return nil, fmt.Errorf("Directory Backed Stash missing inventory Directory")
	}
	return nil, nil
}

func NewDirStash(dirname string) (*DirStash) {
	return &DirStash{InventoryDir: dirname}
}

func (s *DirStash) Load() error {
	// TODO
	return nil
}

func (s *DirStash) IsLoaded() bool {
	// TODO
	return false
}

func (s *DirStash) CreateNew() error {
	// Check if the inventory directory is there
	f, err := os.Open(s.InventoryDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Unable to find inventory dir %s: %s", s.InventoryDir, err)
	}
	if os.IsNotExist(err) {
		// If it isn't there, attempt to create it
		err = os.Mkdir(s.InventoryDir, 0755)
		if err != nil {
			return fmt.Errorf("Unable to make inventory dir %s: %s", s.InventoryDir, err)
		}
	} else {
		// Is it is check to see that it's a directory and empty
		fi, err := f.Stat()
		if err != nil {
			return fmt.Errorf("Unable to validate inventory dir %s: %s", s.InventoryDir, err)
		}
		if !fi.IsDir() {
			return fmt.Errorf("inventory dir is not a directory %s", s.InventoryDir)
		}
		fis, err := f.Readdir(3)
		if err != nil {
			return fmt.Errorf("Unable to validate inventory dir %s: %s", s.InventoryDir, err)
		}
		for _, fi := range fis {
			if fi.Name() != "." && fi.Name() != ".." {
				return fmt.Errorf("inventory dir is not empty %s: %s", s.InventoryDir, fi.Name())
			}
		}
		
	}

	// Now, we know we have a directory we can work with, and it's empty
	return nil
}

func (s *DirStash) CreateArchive(filename string) Archive {
	// TODO
	return nil
}

func (s *DirStash) GetArchive(archivename string) Archive {
	// TODO
	return nil
}

func (s *DirStash) GetArchivelist() []string {
	// TODO
	return []string{}
}

func (s *DirStash) GetFileSize(filename string) int64 {
	// TODO
	return 0
}

func (s *DirStash) IsFileExist(sf *util.SnapFile) bool {
	// TODO
	return false
}

func (s *DirStash) Save() error {
	return nil
}

