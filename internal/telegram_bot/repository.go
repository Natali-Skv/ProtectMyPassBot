package telegram_bot

import (
	"errors"
)

type TgBotRepository interface {
	GetUserID(tgID int) (uint64, error)
}

var (
	NoSuchUserErr    = errors.New("no such user")
	GettingUserIDErr = errors.New("error getting userID by tg-ID")
)
