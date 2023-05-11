package telegram_bot

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"log"
)

type TelegramBot struct {
	l *zap.Logger
}

const (
	setCommand = "set"
	delCommand = "del"
	getCommand = "get"
)

var (
	InitializeError = errors.New("initialize new telegram bot error")
)

type TelegramBotConfig struct {
	Token   string `yaml:"token"`
	Timeout int    `yaml:"timeout"`
}

func NewTelegramBot(logger *zap.Logger) *TelegramBot {
	return &TelegramBot{l: logger}
}

func (tg *TelegramBot) Run(token string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return errors.Join(InitializeError, err)
	}

	tg.l.Info("tg bot authorized", zap.String("tgBotName:", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	//TODO:надо подумать как тут лучше, мб в отдельной горутине
	for update := range updates {
		if update.Message == nil { // ignore non-messages
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		case setCommand:
			//reply := handleSetCommand(update.Message, client)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ok")
			bot.Send(msg)
		case getCommand:
			//reply := handleGetCommand(update.Message, client)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ok")
			bot.Send(msg)
		case delCommand:
			//reply := handleDelCommand(update.Message, client)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ok")
			bot.Send(msg)
		default:
			fmt.Println(update.Message.Command())
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command")
			bot.Send(msg)
		}
	}
	return nil
}
