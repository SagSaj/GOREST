package config

import (
	"github.com/json-iterator/go"
	"io/ioutil"
	"path/filepath"
)

type StConfig struct {
	BindPorts []string
}

func Config_init(path string) *StConfig {
	dir, _ := filepath.Abs("./")
	raw, err := ioutil.ReadFile(dir + path)
	if err != nil {
		panic(err)
	}

	var Conf StConfig
	if err = jsoniter.ConfigFastest.Unmarshal(raw, &Conf); err != nil {
		panic(err)
	}
	return &Conf
}
