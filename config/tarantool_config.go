package config

type TarantoolConfig struct {
	Address       string `yaml:"addr"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	Timeout       int    `yaml:"timeout"`
	Reconnect     int    `yaml:"reconnect"`
	MaxReconnects uint   `yaml:"maxReconnects"`
}
