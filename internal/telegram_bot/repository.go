package telegram_bot

import (
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
)

type TgBotRepository interface {
	GetUserID(tgID models.TgUserID) (models.UserID, error)
	RegisterUser(tgID models.TgUserID, userID models.UserID) error
}
