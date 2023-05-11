package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

func ReadConfig(configPath string, dst interface{}) error {
	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)

	d := yaml.NewDecoder(configFile)
	if err := d.Decode(dst); err != nil {
		return err
	}
	return nil
}
