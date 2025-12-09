package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadYamlConfig(path string, config any) error {
	var err error

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}
