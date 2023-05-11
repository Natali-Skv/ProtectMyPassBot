package main

import (
	"flag"
	"fmt"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const (
	defaultConfigPath = "config/config.yaml"
)

type Config struct {
	Bot telegram_bot.TelegramBotConfig `yaml:"tgBot"`
}

func main() {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal("zap logger build error")
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(logger)

	configPath, err := ParseConfigPathFlag()
	if err != nil {
		logger.Fatal("bad command line arguments passed", zap.Error(err))
	}

	config, err := ReadConfig(configPath)
	if err != nil {
		logger.Fatal("reading config error", zap.Error(err))
	}

	logger.Debug("config", zap.Any("config", config))

	tgBot := telegram_bot.NewTelegramBot(logger)
	err = tgBot.Run(config.Bot.Token)
	if err != nil {
		logger.Fatal("running telegram bot error", zap.Error(err))
	}
}

func ParseConfigPathFlag() (configPath string, err error) {
	flag.StringVar(&configPath, "config", defaultConfigPath, "path to config file")
	flag.Parse()

	s, err := os.Stat(configPath)
	if err != nil {
		return "", err
	}
	if s.IsDir() {
		return "", fmt.Errorf("'%s' is a directory, not a normal file", configPath)
	}

	return configPath, nil
}

func ReadConfig(configPath string) (config *Config, err error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	d := yaml.NewDecoder(configFile)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

//
//func handleSetCommand(message *tgbotapi.Message, client *tarantool.Connection) string {
//	args := message.CommandArguments()
//	if len(args) == 0 {
//		return "Usage: /set <service> <login> <password>"
//	}
//
//	parts := strings.Split(args, " ")
//	if len(parts) != 3 {
//		return "Usage: /set <service> <login> <password>"
//	}
//
//	service := parts[0]
//	login := parts[1]
//	password := parts[2]
//
//	_, err := client.Insert("passwords", []interface{}{message.From.ID, service, login, password})
//	if err != nil {
//		return fmt.Sprintf("Error setting password: %v", err)
//	}
//
//	return "Password saved"
//}
//
//func handleGetCommand(message *tgbotapi.Message, client *tarantool.Connection) string {
//	args := message.CommandArguments()
//	if len(args) == 0 {
//		return "Usage: /get <service>"
//	}
//
//	service := args
//
//	res, err := client.Eval("return box.space.passwords.index.secondary:select({?, ?})", []interface{}{message.From.ID, service})
//	if err != nil {
//		return fmt.Sprintf("Error getting password: %v", err)
//	}
//
//	if len(res.Data) == 0 {
//		return "Password not found"
//	}
//
//	tuple, ok := res.Data[0].([]interface{})
//	if !ok {
//		return "Error getting password"
//	}
//
//	login, ok := tuple[2].(string)
//	if !ok {
//		return "Error getting password"
//	}
//
//	password, ok := tuple[3].(string)
//	if !ok {
//		return "Error getting password"
//	}
//
//	return fmt.Sprintf("Login: %s\nPassword: %s", login, password)
//}
//
//func handleDelCommand(message *tgbotapi.Message, client *tarantool.Connection) string {
//	args := message.CommandArguments()
//	if len(args) == 0 {
//		return "Usage: /del <service>"
//	}
//
//	service := args
//
//	res, err := client.Eval("return box.space.passwords.index.secondary:select({?,?})", []interface{}{message.From.ID, service})
//	if err != nil {
//		return fmt.Sprintf("Error deleting password: %v", err)
//	}
//
//	if len(res.Data) == 0 {
//		return "Password not found"
//	}
//
//	tuple := res.Data[0].([]interface{})
//	_, err = client.Delete("passwords", "primary", tuple[0])
//	if err != nil {
//		return fmt.Sprintf("Error deleting password: %v", err)
//	}
//
//	return "Password deleted"
//}
//
////func handleDelCommand(message *tgbotapi.Message, client *tarantool.Connection) string {
////	args := message.CommandArguments()
////	if len(args) == 0 {
////		return "Usage: /del <service>"
////	}
////
////	service := args
////
////	// Select the tuple to delete
////	resp, err := client.Select("passwords", "secondary", 0, 1, tarantool.IterEq, []interface{}{message.From.ID, service})
////	if err != nil {
////		return fmt.Sprintf("Error deleting password: %v", err)
////	}
////
////	if len(resp.Tuples()) == 0 {
////		return "Password not found"
////	}
////
////	tuple := resp.Tuples()[0]
////
////	// Delete the tuple
////	_, err = client.Delete("passwords", tuple[0])
////	if err != nil {
////		return fmt.Sprintf("Error deleting password: %v", err)
////	}
////
////	return "Password deleted"
////}
//
//func deleteMessageAfterDelay(bot *tgbotapi.BotAPI, message *tgbotapi.Message, delay time.Duration) {
//	time.Sleep(delay)
//
//	deleteMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
//	_, err := bot.DeleteMessage(deleteMsg)
//	if err != nil {
//		log.Printf("Error deleting message: %v", err)
//	}
//}
