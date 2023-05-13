package passman

import (
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
)

type PassmanUsecase interface {
	Get(req GetReqU) (GetRespU, error)
}

type GetReqU struct {
	UserID  models.UserID
	Service string
}

type GetRespU struct {
	Service  string
	Login    string
	Password string
}

var (
	NoSuchUserError = errors.New("No such user error")
)
