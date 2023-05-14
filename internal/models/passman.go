package models

import "errors"

var EmptyUserID UserID = 0
var InvalidUserIDErr = errors.New("invalid user id")

var PassmanRepoErrors = struct {
	NoSuchUserOrServiceInDBErr    error
	CallingGetUserCredesDBFuncErr error

	GettingNextSequenceUserIDErr   error
	InsertingNewUserToUserCredsErr error
	AddUserCredsError              error
	NoSuchUserErr                  error

	DeleteUserCredsErr error
	NoSuchServiceErr   error
}{
	NoSuchUserOrServiceInDBErr:    errors.New("no such user or user service in db error"),
	CallingGetUserCredesDBFuncErr: errors.New("error calling tarantool getUserCredesFunc error"),

	GettingNextSequenceUserIDErr:   errors.New("error getting next sequence userID error"),
	InsertingNewUserToUserCredsErr: errors.New("error inserting new user to user_credentials error"),
	AddUserCredsError:              errors.New("error adding user service credentials to db"),
	NoSuchUserErr:                  errors.New("no such user error"),

	DeleteUserCredsErr: errors.New("error deleting user service credentials"),
	NoSuchServiceErr:   errors.New("no such service credentials"),
}

var PassmanUsecaseErrors = struct {
	UnknownGettingUserCredsErr error
	NoSuchUserOrServiceErr     error
}{
	UnknownGettingUserCredsErr: errors.New("unknown error getting user service credentials"),
	NoSuchUserOrServiceErr:     errors.New("no such user or user service error"),
}

type GrpcError struct {
	Error error
	Code  int
}

var PassmanHandlerErrors = struct {
	UnknownGettingUserCredsErr GrpcError
	NoSuchUserOrServiceErr     GrpcError
}{
	UnknownGettingUserCredsErr: GrpcError{Error: errors.New("unknown error getting user service credentials"), Code: 1},
	NoSuchUserOrServiceErr:     GrpcError{Error: errors.New("no such user or user service error"), Code: 2},
}

type UserID int

func (id UserID) Int() int {
	return int(id)
}

func (id UserID) Int64() int64 {
	return int64(id)
}

func (id UserID) IsValid() bool {
	return id > 0
}

type GetReqR struct {
	UserID  UserID
	Service string
}

type GetRespR struct {
	Service  string
	Login    string
	Password string
}

type GetReqU struct {
	UserID  UserID
	Service string
}

type GetRespU struct {
	Service  string
	Login    string
	Password string
}

type AddCredsData struct {
	Service  string
	Login    string
	Password string
}

type AddCredsReqR struct {
	UserID UserID
	Data   AddCredsData
}

type DeleteCredsReqR struct {
	UserID  UserID
	Service string
}
