package models

import "errors"

type TgUserID int

var InvalidTgUserIDErr = errors.New("invalid telegram id")

func (id TgUserID) Int() int {
	return int(id)
}

func (id TgUserID) Int64() int64 {
	return int64(id)
}

func (id TgUserID) IsValid() bool {
	return id > 0
}

type GetCommandReqU struct {
	TgID       TgUserID
	ArgsString string
}

type GetCommandRespU struct {
	Login    string
	Password string
	Service  string
}

var TgBotUsecaseErrors = struct {
	WrongArgCountErr       error
	GetingUserIDUnknownErr error
	NoSuchUserErr          error
	GettingUserCredsErr    error
	NoSuchCredsErr         error
}{
	WrongArgCountErr:       errors.New("wrong number of command arguments"),
	GetingUserIDUnknownErr: errors.New("unknown error getting userID"),
	NoSuchUserErr:          errors.New("no such user"),
	GettingUserCredsErr:    errors.New("error getting user creds with passman"),
	NoSuchCredsErr:         errors.New("no credentials for such service"),
}

var TgBotRepoErrors = struct {
	NoSuchUserErr    error
	GettingUserIDErr error
	RegisterUserErr  error
}{
	NoSuchUserErr:    errors.New("no such user"),
	GettingUserIDErr: errors.New("error getting userID by tg-ID"),
	RegisterUserErr:  errors.New("error registering user in tg_id space"),
}
