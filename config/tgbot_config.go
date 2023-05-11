package config

const (
	TgBotDefaultConfigPath = "config/yaml/tgbot.yaml"
)

type TelegramBotConfig struct {
	Token   string `yaml:"token"`
	Timeout int    `yaml:"timeout"`
}

type BotConfig struct {
	Bot         TelegramBotConfig `yaml:"tgBot"`
	PassmanAddr string            `yaml:"passmanAddr"`
}
