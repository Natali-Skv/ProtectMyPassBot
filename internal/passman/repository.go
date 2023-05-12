package passman

import "errors"

type PassmanRepository interface {
	Get(req GetReqR) (GetRespR, error)
}

type GetReqR struct {
	UserID  uint64
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
