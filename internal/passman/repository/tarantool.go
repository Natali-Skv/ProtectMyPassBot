package repository

import (
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/passman"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
)

const (
	userCredsSpace    = "user_credentials"
	getUserCredesFunc = "get_user_creds"
	primary           = "primary"
)

type PassmanRepo struct {
	l    *zap.Logger // TODO он точно нужен? или все логируется выше
	conn *tarantool.Connection
}

func NewPassmanRepo(l *zap.Logger, conn *tarantool.Connection) *PassmanRepo {
	if l == nil || conn == nil {
		return nil
	}
	return &PassmanRepo{l: l, conn: conn}
}

func (pr *PassmanRepo) Get(req passman.GetReqR) (passman.GetRespR, error) {
	var queryResult []struct {
		//Service  string
		Login    string
		Password string
	}

	pr.l.Debug("req", zap.Any("user_id", req.UserID), zap.String("servicename", req.Service))
	err := pr.conn.CallTyped(getUserCredesFunc, []interface{}{req.UserID, req.Service}, &queryResult)
	pr.l.Debug("result", zap.Any("res", queryResult))

	if err != nil {
		pr.l.Error("error quering user credentials", zap.Error(err))
		return passman.GetRespR{}, errors.Join(telegram_bot.GettingUserIDErr, err)
	}
	return passman.GetRespR{
		Service:  req.Service,
		Login:    queryResult[0].Login,
		Password: queryResult[0].Password,
	}, nil
}
