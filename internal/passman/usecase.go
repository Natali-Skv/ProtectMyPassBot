package passman

import "errors"

type PassmanUsecase interface {
	Get(req GetReqU) (GetRespU, error)
}

type GetReqU struct {
	UserID  uint64
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
