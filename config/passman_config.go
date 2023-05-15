package config

const (
	PassmanDefaultConfigPath = "/home/passman.yaml"
)

type PassmanConfig struct {
	BindAddr  string          `yaml:"bindAddr"`
	Tarantool TarantoolConfig `yaml:"tarantool"`
}
