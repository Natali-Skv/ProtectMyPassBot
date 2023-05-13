package telegram_bot

import (
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
)

type TgBotUsecase interface {
	GetCommand(req *models.GetCommandReqU) (resp *models.GetCommandRespU, err error)
	RegisterUser(tgID models.TgUserID) (userID models.UserID, err error)
}
