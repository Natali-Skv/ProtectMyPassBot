package passman

import (
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
)

type PassmanRepository interface {
	Get(req GetReqR) (GetRespR, error)
}

type GetReqR struct {
	UserID  models.UserID
	Service string
}

type GetRespR struct {
	Service  string
	Login    string
	Password string
}

var (
	NoSuchUserInDBError = errors.New("No such user error")
)
