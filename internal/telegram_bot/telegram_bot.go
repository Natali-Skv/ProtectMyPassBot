package telegram_bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/Natali-Skv/ProtectMyPassBot/config"
	passmanProto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"log"
)

const (
	setCommand = "set"
	delCommand = "del"
	getCommand = "get"
)

var (
	InitializeError = errors.New("initialize new telegram bot error")
)

type TelegramBot struct {
	l          *zap.Logger
	passmanCli passmanProto.PassmanServiceClient
	config     config.TelegramBotConfig
	r          TgBotRepository
}

func NewTelegramBot(tgbotConfig config.TelegramBotConfig, logger *zap.Logger, passmanCli passmanProto.PassmanServiceClient, r TgBotRepository) *TelegramBot {
	return &TelegramBot{
		l:          logger,
		passmanCli: passmanCli,
		config:     tgbotConfig,
		r:          r,
	}
}

func (b *TelegramBot) Run() error {
	bot, err := tgbotapi.NewBotAPI(b.config.Token)
	if err != nil {
		return errors.Join(InitializeError, err)
	}

	b.l.Info("b bot authorized", zap.String("tgBotName:", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.config.Timeout

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
			msg := tgbotapi.MessageConfig{}

			b.l.Debug(update.Message.From.UserName)
			b.l.Debug("TG ID", zap.Int("ID", update.Message.From.ID))

			userID, err := b.r.GetUserID(update.Message.From.ID)

			if err != nil {
				b.l.Error("error getting userID by tgID", zap.Error(err))
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "error")
				bot.Send(msg)
				break
				//	TODO возврат сообщения об ошибке
			}
			credentials, err := b.passmanCli.GetCredentials(context.TODO(), &passmanProto.GetReq{UserID: userID, ServiceName: "VK"})
			b.l.Info("passman resp", zap.Any("creds", credentials), zap.Error(err))
			// TODOсделать пароль и логин удобными для копирования
			msgText := fmt.Sprintf("Service: `%s` \nLogin: `%s`\nPassword: `%s`", credentials.ServiceName, credentials.Login, credentials.Password)
			//credentials, err := b.passmanCli.GetCredentials(context.Background(), &passmanProto.GetReq{
			//	Token:       "TODO change to userID",
			//	ServiceName: "VK",
			//})
			//if err != nil {
			//	b.l.Error("error getting credentials", zap.Error(err))
			//	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "error")
			//}

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			msg.ParseMode = "Markdown"
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
