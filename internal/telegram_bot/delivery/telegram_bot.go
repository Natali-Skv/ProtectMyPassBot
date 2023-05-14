package delivery

import (
	"context"
	"errors"
	"fmt"
	"github.com/Natali-Skv/ProtectMyPassBot/config"
	m "github.com/Natali-Skv/ProtectMyPassBot/internal/models"
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
	u      m.TgBotUsecase
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

func NewTelegramBot(tgbotConfig config.TelegramBotConfig, logger *zap.Logger, u m.TgBotUsecase, dtm delay_task_manager.DelayTaskManager) *TelegramBotHandler {
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
		case tgbot.HelpCommand.Name, tgbot.StartCommand.Name:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, tgbot.HelpCommand.RespFmtString)
			msg.ParseMode = MarkdownMode
			send, err := bot.Send(msg)
			if err != nil {
				tgb.l.Error("error sending message", zap.Error(err), zap.Any("send", send))
			}
		case tgbot.SetCommand.Name:
			msg := tgb.setCommand(&update)
			send, err := bot.Send(msg)
			if err != nil {
				tgb.l.Error("error sending message", zap.Error(err), zap.Any("send", send))
			}
			tgb.dtm.AddTask(&messageDeleter{
				deleteMsg: tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID),
				l:         tgb.l,
				bot:       bot,
			})
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
	respU, err := tgb.u.GetCommand(&m.GetCommandReqU{TgID: m.TgUserID(update.Message.From.ID), ArgsString: update.Message.CommandArguments()})
	switch {
	case err == nil:
		ok = true
		msg.Text = fmt.Sprintf(tgbot.GetCommand.RespFmtString, respU.Service, respU.Login, respU.Password)
	case errors.Is(err, m.TgBotUsecaseErrors.WrongArgCountErr):
		msg.Text = tgbot.GetCommand.Usage
	case errors.Is(err, m.TgBotUsecaseErrors.NoSuchCredsErr):
		msg.Text = tgbot.NoSuchCredsMsg
	case errors.Is(err, m.TgBotUsecaseErrors.NoSuchUserErr):
		msg.Text = tgbot.NoSuchUserMsg
	default:
		tgb.l.Error("error handling command /get", zap.Error(err), zap.Any("upd", update))
		msg.Text = tgbot.UnknownErrorResp
	}
	msg.ParseMode = MarkdownMode
	return msg, ok
}

func (tgb *TelegramBotHandler) setCommand(update *tgbotapi.Update) (msg tgbotapi.MessageConfig) {
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
	serviceName, err := tgb.u.SetCommand(&m.SetCommandReqU{TgID: m.TgUserID(update.Message.From.ID), ArgsString: update.Message.CommandArguments()})
	switch {
	case err == nil:
		msg.Text = fmt.Sprintf(tgbot.SetCommand.RespFmtString, serviceName)
	case errors.Is(err, m.TgBotUsecaseErrors.WrongArgCountErr):
		msg.Text = tgbot.SetCommand.Usage
	default:
		tgb.l.Error("error handling command /set", zap.Error(err), zap.Any("upd", update))
		msg.Text = tgbot.UnknownErrorResp
	}
	msg.ParseMode = MarkdownMode
	return msg
}
