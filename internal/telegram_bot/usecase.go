package telegram_bot

import (
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
)

type GetCommandReqU struct {
	TgID       models.TgUserID
	ArgsString string
}

type TgBotUsecase interface {
	GetCommand(req GetCommandReqU) (resp string, err error)
	RegisterUser(tgID models.TgUserID) (userID models.UserID, err error)
}

var (
	WrongArgCountErr = errors.New("wrong number of command arguments")
)
