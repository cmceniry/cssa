package command

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"fmt"
	"os"
	"github.com/cmceniry/cssa/stash"
)

type GeneralOptions struct {
	CssaConfigfile	string
	ConfigFile	string
	CstarConfigfile	string
	CassandraConfigFile	string
	Stash		stash.Stash
}

var Commands = map[string]func(*GeneralOptions, []string){
	"archive": Archive,
	"init": InitializeStash,
	"list": List,
	"size": Size,
	"sizesummary": SizeSummary,
}

func ParseConfigFile(filename string) (*GeneralOptions, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	d := make(map[string]interface{})
	yaml.Unmarshal(data, &d)

	ret := &GeneralOptions{
		CssaConfigfile: filename,
		ConfigFile: filename,
	}
	if val, ok := d["cassandraconfig"]; ok {
		if v, ok := val.(string); ok {
			ret.CstarConfigfile = v
			ret.CassandraConfigFile = v
		} else {
			return nil, fmt.Errorf("Invalid Config: cassandraconfig: %s(%T)\n", val, val)
		}
	}
	if val, ok := d["stash"]; ok {
		if v, ok := val.(map[interface{}]interface{}); ok {
			vstring := make(map[string]interface{})
			for keyintf, valintf := range v {
				if key, ok := keyintf.(string); ok {
					vstring[key] = valintf
				} else {
					return nil, fmt.Errorf("Invalid stash key: %s should be string", keyintf)
				}
			}
			s, err := stash.ParseConfig(vstring)
			if err != nil {
				return nil, fmt.Errorf("Invalid Config: stash: %s", err)
			}
			ret.Stash = s
		} else {
			return nil, fmt.Errorf("Invalid Config: stash: %s(%T)\n", val, val)
		}
	}


	return ret, nil
}
