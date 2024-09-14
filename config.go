package main

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type SyncConfig struct {
	CloudDir   string   `yaml:"cloud_dir"`
	KeepLatest bool     `yaml:"keep_latest"`
	Patterns   []string `yaml:"patterns"`
}

var syncConfig SyncConfig

func loadSyncConfig() error {
	data, err := ioutil.ReadFile("sync.yaml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &syncConfig)
}
