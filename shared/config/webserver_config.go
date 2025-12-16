package config

type WebServerConfig struct {
	Port uint16 `yaml:"port"`

	ServerName string `yaml:"server_name"`

	MongoUri       string `yaml:"mongo_uri"`
	MongoEnvDBName string `yaml:"mongo_env_db_name"`
}
