package config

type TarantoolConfig struct {
	Address  string `yaml:"addr"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
