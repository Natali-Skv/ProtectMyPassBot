package passman

import "errors"

type PassmanUsecase interface {
	Get(req GetReqU) (GetRespU, error)
}

type UserID int

type GetReqU struct {
	UserID  UserID
	Service string
}

type Credentials struct {
	Login    string
	Password string
}

type GetRespU struct {
	Service     string
	Credentials []Credentials
}

var (
	NoSuchUserError = errors.New("No such user error")
)
