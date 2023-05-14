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
	NoSuchServiceErr       error
	GettingUserCredsErr    error
	NoSuchCredsErr         error
	RegisterUserErr        error
	SetUserCredsErr        error
	DelUserCredsErr        error
}{
	WrongArgCountErr:       errors.New("wrong number of command arguments"),
	GetingUserIDUnknownErr: errors.New("unknown error getting userID"),
	NoSuchUserErr:          errors.New("no such user"),
	NoSuchServiceErr:       errors.New("no such user service"),
	GettingUserCredsErr:    errors.New("error getting user creds with passman"),
	NoSuchCredsErr:         errors.New("no credentials for such service"),
	RegisterUserErr:        errors.New("error registering user"),
	SetUserCredsErr:        errors.New("error setting credentials for user"),
	DelUserCredsErr:        errors.New("error deleting credentials for user"),
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

type SetCommandReqU struct {
	TgID       TgUserID
	ArgsString string
}

type DelCommandReqU struct {
	TgID       TgUserID
	ArgsString string
}

type TgBotUsecase interface {
	GetCommand(req *GetCommandReqU) (resp *GetCommandRespU, err error)
	SetCommand(req *SetCommandReqU) (serviceName string, err error)
	DelCommand(req *DelCommandReqU) (serviceName string, err error)
}

type TgBotRepository interface {
	GetUserID(tgID TgUserID) (UserID, error)
	RegisterUser(tgID TgUserID, userID UserID) error
}
