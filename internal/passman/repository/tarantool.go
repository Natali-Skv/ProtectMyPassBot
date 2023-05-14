package repository

import (
	"errors"
	m "github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
)

const (
	userCredsSpace    = "user_credentials"
	getUserCredesFunc = "get_user_creds"
	primary           = "primary"
	addUser           = "add_user"
	addUserService    = "add_user_service"
	deleteUserService = "remove_user_service"

	addUserServiceCodeOK         = 0
	addUserServiceCodeNoSuchUser = 1

	deleteUserServiceCodeOK            = 0
	deleteUserServiceCodeNoSuchService = 1
	deleteUserServiceCodeNoSuchUser    = 2
)

type UserCredsTuple struct {
	Login    string
	Password string
}

func (u UserCredsTuple) isEmpty() bool {
	return u.Login == "" && u.Password == ""
}

type PassmanRepo struct {
	l    *zap.Logger
	conn *tarantool.Connection
}

func NewPassmanRepo(l *zap.Logger, conn *tarantool.Connection) *PassmanRepo {
	if l == nil || conn == nil {
		return nil
	}
	return &PassmanRepo{l: l, conn: conn}
}

func (pr *PassmanRepo) Get(req m.GetReqR) (m.GetRespR, error) {
	var queryResult []UserCredsTuple

	pr.l.Debug("req", zap.Any("user_id", req.UserID), zap.String("servicename", req.Service))
	err := pr.conn.CallTyped(getUserCredesFunc, []interface{}{req.UserID, req.Service}, &queryResult)
	pr.l.Debug("result", zap.Any("res", queryResult))

	if err != nil {
		return m.GetRespR{}, errors.Join(m.PassmanRepoErrors.CallingGetUserCredesDBFuncErr, err)
	}

	if queryResult[0].isEmpty() {
		return m.GetRespR{}, m.PassmanRepoErrors.NoSuchUserOrServiceInDBErr
	}

	return m.GetRespR{
		Service:  req.Service,
		Login:    queryResult[0].Login,
		Password: queryResult[0].Password,
	}, nil
}

func (pr *PassmanRepo) Register() (m.UserID, error) {
	var userIDSequenceResp []struct {
		UserID m.UserID
	}
	err := pr.conn.CallTyped(addUser, []interface{}{map[interface{}]interface{}{}}, &userIDSequenceResp)
	if err != nil {
		pr.l.Debug("error inserting next userIDSequenceResp", zap.Error(err))
		return m.EmptyUserID, errors.Join(m.PassmanRepoErrors.AddNewUserToUserCredsErr, err)
	}

	return userIDSequenceResp[0].UserID, nil
}

func (pr *PassmanRepo) AddCredentials(req m.AddCredsReqR) error {
	var addCredsResp []struct {
		Error string
		Code  int64
	}
	err := pr.conn.CallTyped(addUserService, []interface{}{req.UserID, req.Data.Service, req.Data.Login, req.Data.Password}, &addCredsResp)
	if err != nil {
		pr.l.Debug("error adding user creds", zap.Error(err), zap.Any("resp", addCredsResp))
		return errors.Join(m.PassmanRepoErrors.AddUserCredsError, err)
	}
	switch addCredsResp[0].Code {
	case addUserServiceCodeOK:
		return nil
	case addUserServiceCodeNoSuchUser:
		return errors.Join(m.PassmanRepoErrors.NoSuchUserErr, errors.New(addCredsResp[0].Error))
	default:
		pr.l.Debug("unknown code in return value by add user creds tarantool func", zap.Error(err), zap.Any("resp", addCredsResp))
		return errors.Join(m.PassmanRepoErrors.AddUserCredsError, errors.New(addCredsResp[0].Error))
	}
}

func (pr *PassmanRepo) DeleteCreds(req m.DeleteCredsReqR) error {
	var delCredsResp []struct {
		Error string
		Code  int64
	}
	err := pr.conn.CallTyped(deleteUserService, []interface{}{req.UserID, req.Service}, &delCredsResp)
	if err != nil {
		pr.l.Debug("error deleting user creds ", zap.Error(err), zap.Any("resp", delCredsResp))
		return errors.Join(m.PassmanRepoErrors.DeleteUserCredsErr, err)
	}
	switch delCredsResp[0].Code {
	case addUserServiceCodeOK:
		return nil
	case deleteUserServiceCodeNoSuchService:
		return errors.Join(m.PassmanRepoErrors.NoSuchServiceErr, errors.New(delCredsResp[0].Error))
	case deleteUserServiceCodeNoSuchUser:
		return errors.Join(m.PassmanRepoErrors.NoSuchUserErr, errors.New(delCredsResp[0].Error))
	default:
		pr.l.Debug("unknown code in return value by delete user creds tarantool func", zap.Error(err), zap.Any("resp", delCredsResp))
		return errors.Join(m.PassmanRepoErrors.DeleteUserCredsErr, errors.New(delCredsResp[0].Error))
	}
}
