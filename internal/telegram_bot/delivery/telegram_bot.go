package delivery

import (
	"context"
	"errors"
	"fmt"
	"github.com/Natali-Skv/ProtectMyPassBot/config"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	tgbot "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/tools/delay_task_manager"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"time"
)

var (
	InitializeError = errors.New("initialize new telegram bot error")
)

const (
	MarkdownMode = "Markdown"
)

type TelegramBotHandler struct {
	l      *zap.Logger
	u      tgbot.TgBotUsecase
	config config.TelegramBotConfig
	dtm    delay_task_manager.DelayTaskManager
}

type messageDeleter struct {
	deleteMsg tgbotapi.DeleteMessageConfig
	l         *zap.Logger
	bot       *tgbotapi.BotAPI
}

func (md *messageDeleter) Process() {
	if _, err := md.bot.Send(md.deleteMsg); err != nil {
		md.l.Error("error sending delete-message by tg-bot", zap.Error(err))
	}
}

func NewTelegramBot(tgbotConfig config.TelegramBotConfig, logger *zap.Logger, u tgbot.TgBotUsecase, dtm delay_task_manager.DelayTaskManager) *TelegramBotHandler {
	go dtm.Run(context.Background(), time.Second*time.Duration(tgbotConfig.GetCommandDeleteTimout))
	return &TelegramBotHandler{
		l:      logger,
		config: tgbotConfig,
		u:      u,
		dtm:    dtm,
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
		case tgbot.SetCommand.Name:
			//reply := handleSetCommand(update.Message, client)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ok")
			bot.Send(msg)
		case tgbot.GetCommand.Name:
			msg, ok := tgb.getCommand(&update)
			send, err := bot.Send(msg)
			if err != nil {
				tgb.l.Error("error sending message", zap.Error(err), zap.Any("send", send))
			}
			if ok {
				tgb.dtm.AddTask(&messageDeleter{
					deleteMsg: tgbotapi.NewDeleteMessage(send.Chat.ID, send.MessageID),
					l:         tgb.l,
					bot:       bot,
				})
			}

		case tgbot.DelCommand.Name:
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

func (tgb *TelegramBotHandler) getCommand(update *tgbotapi.Update) (msg tgbotapi.MessageConfig, ok bool) {
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
	var err error
	respU, err := tgb.u.GetCommand(&models.GetCommandReqU{TgID: models.TgUserID(update.Message.From.ID), ArgsString: update.Message.CommandArguments()})
	switch {
	case err == nil:
		ok = true
		msg.Text = fmt.Sprintf(tgbot.GetCommand.RespFmtString, respU.Service, respU.Login, respU.Password)
	case errors.Is(err, models.TgBotUsecaseErrors.WrongArgCountErr):
		msg.Text = tgbot.GetCommand.Usage
	case errors.Is(err, models.TgBotUsecaseErrors.NoSuchCredsErr):
		msg.Text = tgbot.NoSuchCredsMsg
	default:
		tgb.l.Error("error handling command /get", zap.Error(err))
		msg.Text = tgbot.UnknownErrorResp
	}
	msg.ParseMode = MarkdownMode
	return msg, ok
}
