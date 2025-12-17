package config

type ClientConfig struct {
	Url  string `yaml:"url"`
	Port uint16 `yaml:"port"`
}
