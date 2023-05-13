package delivery

import (
	"errors"
	"fmt"
	"github.com/Natali-Skv/ProtectMyPassBot/config"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	tgbot "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

var (
	InitializeError = errors.New("initialize new telegram bot error")
)

const (
	MarkdownMode     = "Markdown"
	UnknownErrorResp = "Sorry, unknown error occurred, try again later."
)

type TelegramBotHandler struct {
	l      *zap.Logger
	u      tgbot.TgBotUsecase
	config config.TelegramBotConfig
}

func NewTelegramBot(tgbotConfig config.TelegramBotConfig, logger *zap.Logger, u tgbot.TgBotUsecase) *TelegramBotHandler {
	return &TelegramBotHandler{
		l:      logger,
		config: tgbotConfig,
		u:      u,
	}
}

func (tgb *TelegramBotHandler) Run() error {
	bot, err := tgbotapi.NewBotAPI(tgb.config.Token)
	if err != nil {
		return errors.Join(InitializeError, err)
	}

	tgb.l.Info("tgb bot authorized", zap.String("tgBotName:", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = tgb.config.Timeout

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case models.SetCommand.Name:
			//reply := handleSetCommand(update.Message, client)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ok")
			bot.Send(msg)
		case models.GetCommand.Name:
			msg := tgb.getCommand(&update)
			send, err := bot.Send(msg)
			if err != nil {
				tgb.l.Error("error sending message", zap.Error(err), zap.Any("send", send))
			}
		case models.DelCommand.Name:
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

func (tgb *TelegramBotHandler) getCommand(update *tgbotapi.Update) *tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	var err error
	msg.Text, err = tgb.u.GetCommand(tgbot.GetCommandReqU{TgID: models.TgUserID(update.Message.From.ID), ArgsString: update.Message.CommandArguments()})

	if err != nil {
		switch err {
		case tgbot.WrongArgCountErr:
			msg.Text = models.GetCommand.Usage
		default:
			tgb.l.Error("error handling command /get", zap.Error(err))
			msg.Text = UnknownErrorResp
		}
	}
	msg.ParseMode = MarkdownMode
	return &msg
}
