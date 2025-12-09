package config

type WebServerConfig struct {
	Port uint16 `yaml:"port"`

	MongoUri string `yaml:"mongo_uri"`
}
