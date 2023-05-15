package config

const (
	TgBotDefaultConfigPath = "/home/tgbot.yaml"
)

type TelegramBotConfig struct {
	Token                  string `yaml:"token"`
	Timeout                int    `yaml:"timeout"`
	GetCommandDeleteTimout int    `yaml:"getCommandDeleteTimout"`
}

type BotConfig struct {
	Bot         TelegramBotConfig `yaml:"tgBot"`
	PassmanAddr string            `yaml:"passmanAddr"`
	Tarantool   TarantoolConfig   `yaml:"tarantool"`
}
