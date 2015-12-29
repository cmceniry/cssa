package stash

import (
	yaml "gopkg.in/yaml.v2"
	"os"
	"io/ioutil"
	"github.com/cmceniry/cssa/util"
	"fmt"
	"time"
	"path/filepath"
)

var (
	DuplicateFileExists = fmt.Errorf("Duplicate File Exists")
	ConflictFileExists = fmt.Errorf("Conflict File Exists")
)

type HollowFile struct {
	Parent		*HollowStash	`yaml:"-"`
	Name		string
	OriginalPath	string
	Size		int64
	OriginalInode	uint64
}


/*
 *
 * Hollow Stash
 *
 */

type HollowStash struct {
	InventoryFile	string
	Loaded		bool	`yaml:"-"`
	Files		map[string]*HollowFile
	Archives	map[string]*HollowArchive
}

func NewHollowStashFromConfig(c map[string]interface{}) (*HollowStash, error) {
	if filenameintf, ok := c["inventoryfile"]; ok {
		if filename, ok := filenameintf.(string); ok {
			r := NewHollowStash(filename)
			return r, nil
		} else {
			return nil, fmt.Errorf("Invalid inventoryname: %s(%T)", filenameintf, filenameintf)
		}
	} else {
		return nil, fmt.Errorf("Hollow Stash missing inventory file")
	}
}

func NewHollowStash(filename string) (*HollowStash) {
	return &HollowStash{InventoryFile: filename}
}

func (s *HollowStash) Load() error {
	if s.Loaded {
		return fmt.Errorf("Already Loaded")
	}
	f, err := os.Open(s.InventoryFile)
	if err != nil {
		return err
	}
	defer f.Close()
	y, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	r := &HollowStash{}
	err = yaml.Unmarshal(y, r)
	if err != nil {
		return err
	}
	if r.Files != nil {
		s.Files = r.Files
	} else {
		s.Files = make(map[string]*HollowFile)
	}
	if r.Archives != nil {
		s.Archives = r.Archives
	} else {
		s.Archives = make(map[string]*HollowArchive)
	}
	for _, f := range s.Files {
		f.Parent = s
	}
	for _, a := range s.Archives {
		a.Parent = s
	}
	s.Loaded = true
	return nil
}

func (s *HollowStash) IsLoaded() bool {
	return s.Loaded
}

func (s *HollowStash) CreateNew() error {
	m, err := yaml.Marshal(*s)
	if err != nil {
		return err
	}
	f, err := os.Open(s.InventoryFile)
	if os.IsNotExist(err) {
		f, err = os.Create(s.InventoryFile)
		if err != nil {
			return err
		}
		f.Write(m)
		f.Close()
	} else if err != nil {
		return err
	} else {
		f.Close()
		return os.ErrExist
	}
	return nil
}

func (s *HollowStash) AddFile(newfile *HollowFile) error {
	if oldfile, ok := s.Files[newfile.Name]; ok {
		if oldfile.Name == newfile.Name && oldfile.OriginalInode == newfile.OriginalInode && oldfile.Size == newfile.Size {
			return DuplicateFileExists
		} else {
			return ConflictFileExists
		}
	}
	s.Files[newfile.Name] = newfile
	return nil
}

func (s *HollowStash) CreateArchive(archivename string) Archive {
	return &HollowArchive{Parent: s, Starttime: time.Now().Unix(), Name: archivename, Files: []string{}}
}

func (s *HollowStash) GetArchive(name string) Archive {
	if ret, ok := s.Archives[name]; ok {
		return ret
	} else {
		return nil
	}
}

func (s *HollowStash) GetArchivelist() []string {
	ret := []string{}
	for a, _ := range s.Archives {
		ret = append(ret, a)
	}
	return ret
}

func (s *HollowStash) GetFilelist() []string {
	files := []string{}
	for _, f := range s.Files {
		files = append(files, f.Name)
	}
	return files
}

func (s *HollowStash) GetFileSize(filename string) int64 {
	if f, ok := s.Files[filename]; ok {
		return f.Size
	} else {
		return -1
	}
}

func (s *HollowStash) IsFileExist(sf *util.SnapFile) bool {
	_, ok := s.Files[sf.Filename]
	if !ok {
		return false
	}
	return true
}

func (s *HollowStash) Save() (error) {
	dir, filename := filepath.Split(s.InventoryFile)
	f, err := ioutil.TempFile(dir, filename)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := yaml.Marshal(*s)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	err = os.Rename(s.InventoryFile, s.InventoryFile + ".bak")
	if err != nil {
		return err
	}
	os.Rename(f.Name(), s.InventoryFile)
	if err != nil {
		return err
	}
	return nil
}


/*
 *
 * Hollow Archive
 *
 */


type HollowArchive struct {
	Parent		*HollowStash	`yaml:"-"`
	Name		string
	Starttime	int64
	Committime	int64
	Files		[]string
	Invalid		bool
}

func (a *HollowArchive) AddFile(f *util.SnapFile) error {
	hf := &HollowFile{
		Parent: a.Parent,
		Name: f.Filename,
		OriginalPath: f.Filepath,
		Size: f.Size,
		OriginalInode: f.Inode,
	}
	err := a.Parent.AddFile(hf)
	if err != nil && err != DuplicateFileExists {
		return err
	}
	a.Files = append(a.Files, hf.Name)
	return err
}

func (a *HollowArchive) Commit() error {
	a.Parent.Archives[a.Name] = a
	a.Committime = time.Now().Unix()
	err := a.Parent.Save()
	return err
}

func (a *HollowArchive) GetFilelist() []string {
	return a.Files
}

func (a *HollowArchive) GetName() string {
	return a.Name
}

func (s *HollowArchive) IsValid() bool {
	return !s.Invalid
}

func (s *HollowArchive) SetInvalid() {
	s.Invalid = true
}

func (s *HollowArchive) SetValid() {
	s.Invalid = false
}

