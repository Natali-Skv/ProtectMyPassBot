package telegram_bot

import (
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
)

type TgBotRepository interface {
	GetUserID(tgID models.TgUserID) (models.UserID, error)
	RegisterUser(tgID models.TgUserID, userID models.UserID) error
}

var (
	NoSuchUserErr    = errors.New("no such user")
	GettingUserIDErr = errors.New("error getting userID by tg-ID")
	RegisterUserErr  = errors.New("error registering user in tg_id space")
)
