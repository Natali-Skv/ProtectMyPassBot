package models

import "errors"

type Command struct {
	Name          string
	Usage         string
	ArgumentCount int
	RespFmtString string
}

var (
	SetCommand                      = Command{Name: "set", Usage: "Usage: `/set <service> <login> <password>`", ArgumentCount: 3}
	DelCommand                      = Command{Name: "del", Usage: "Usage: `/del <service>`", ArgumentCount: 1}
	GetCommand                      = Command{Name: "get", Usage: "Usage: `/get <service>`", ArgumentCount: 1, RespFmtString: "Service: `%s` \nLogin: `%s`\nPassword: `%s`"}
	GetCommandServiceArgumentNumber = 0
)

type TgUserID int

var EmptyTgUserID TgUserID = 0
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
