package util

import (
	yaml "gopkg.in/yaml.v2"
	"os"
	"syscall"
	"io/ioutil"
	"path/filepath"
	"strings"
	"fmt"
)

var (
	SnapshotNotFound = fmt.Errorf("SnapshotNotFound")
)

func GetDataDirs(configfile string) ([]string, error) {
	ret := []string{}

	f, err := os.Open(configfile)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	d := make(map[string]interface{})
	yaml.Unmarshal(data, &d)
	if val, ok := d["data_file_directories"]; ok {
		if v, ok := val.([]interface{}); ok {
			for _, e := range v {
				if s, ok := e.(string); ok {
					ret = append(ret, s)
				}
			}
		}
		return ret, nil
	} else {
		return []string{}, nil
	}
	return []string{}, nil
}

func GetAllSnapShotDirs(configfile, snapshotname string) ([]string, error) {
	datadirs, err := GetDataDirs("/etc/dse/cassandra/cassandra.yaml")
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, datadir := range datadirs {
		ssds, _ := filepath.Glob(filepath.Join(datadir, "*", "*", "snapshots", snapshotname))
		ret = append(ret, ssds...)
	}
	return ret, nil
}

func GetKsCfFromName(snapshotdirname string) (string, string, bool) {
	r := strings.Split(snapshotdirname, string(filepath.Separator))
	if len(r) < 4 {
		return "", "", false
	}
	return r[len(r)-4], r[len(r)-3], true
}

func GetLocalSnapShot(configfile, snapshotname string) (*SnapManifest, []string, error) {
	ret := &SnapManifest{SnapshotName: snapshotname}
	warns := []string{}
	ssds, err := GetAllSnapShotDirs(configfile, snapshotname)
	if err != nil {
		return nil, nil, err
	}
	if len(ssds) == 0 {
		return nil, nil, SnapshotNotFound
	}
	for _, ssd := range ssds {
		ks, cf, ok := GetKsCfFromName(ssd)
		if !ok {
			warns = appendWarn(warns, "Unable to parse %s", ssd)
			continue
		}
		datafiles, _ := filepath.Glob(filepath.Join(ssd, ks + "-" + cf + "-*.db"))
		for _, datafile := range datafiles {
			_, filename := filepath.Split(datafile)
			fi, err := os.Stat(datafile)
			if err != nil {
				warns = appendWarn(warns, "Unable to stat %s: %s", datafile, err)
				continue
			}
			if !fi.Mode().IsRegular() {
				warns = appendWarn(warns, "%s is not a file", datafile)
				continue
			}
			stat, ok := fi.Sys().(*syscall.Stat_t)
			if !ok {
				warns = appendWarn(warns, "Unable to process stat %s", datafile)
				continue
			}
			sf := SnapFile{Filename: filename, Filepath: datafile, Inode: stat.Ino, Size: fi.Size(), Keyspace: ks, ColumnFamily: cf}
			ret.Files = append(ret.Files, &sf)
		}
	}
	return ret, warns, nil
}
