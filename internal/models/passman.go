package models

import "errors"

type UserID int

func (id UserID) Int() int {
	return int(id)
}

func (id UserID) Int64() int64 {
	return int64(id)
}

var EmptyUserID UserID = 0
var InvalidUserIDErr = errors.New("invalid user id")

func (id UserID) IsValid() bool {
	return id > 0
}
