package config

const (
	PassmanDefaultConfigPath = "config/yaml/passman.yaml"
)

type PassmanConfig struct {
	BindAddr  string          `yaml:"bindAddr"`
	Tarantool TarantoolConfig `yaml:"tarantool"`
}
