package repository

import (
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
)

const (
	userCredsSpace    = "user_credentials"
	getUserCredesFunc = "get_user_creds"
	primary           = "primary"
)

type userCredsTuple struct {
	Login    string
	Password string
}

func (u userCredsTuple) isEmpty() bool {
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

func (pr *PassmanRepo) Get(req models.GetReqR) (models.GetRespR, error) {
	var queryResult []userCredsTuple

	pr.l.Debug("req", zap.Any("user_id", req.UserID), zap.String("servicename", req.Service))
	err := pr.conn.CallTyped(getUserCredesFunc, []interface{}{req.UserID, req.Service}, &queryResult)
	pr.l.Debug("result", zap.Any("res", queryResult))

	if err != nil {
		return models.GetRespR{}, errors.Join(models.PassmanRepoErrors.CallingGetUserCredesDBFuncErr, err)
	}

	if queryResult[0].isEmpty() {
		return models.GetRespR{}, models.PassmanRepoErrors.NoSuchUserOrServiceInDBErr
	}

	return models.GetRespR{
		Service:  req.Service,
		Login:    queryResult[0].Login,
		Password: queryResult[0].Password,
	}, nil
}
