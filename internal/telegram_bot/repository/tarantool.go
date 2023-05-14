package repository

import (
	"errors"
	m "github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
)

const (
	tgIdSpace = "tg_id"
	primary   = "primary"
)

type tgIdTuple struct {
	TgID   m.TgUserID
	UserID m.UserID
}

type TgBotRepo struct {
	l    *zap.Logger
	conn *tarantool.Connection
}

func NewTgBotRepo(l *zap.Logger, conn *tarantool.Connection) *TgBotRepo {
	if l == nil || conn == nil {
		return nil
	}
	return &TgBotRepo{l: l, conn: conn}
}

func (tr *TgBotRepo) GetUserID(tgID m.TgUserID) (userID m.UserID, err error) {
	resultTuple := tgIdTuple{}
	err = tr.conn.GetTyped(tgIdSpace, primary, []interface{}{tgID}, &resultTuple)
	if err != nil {
		return m.EmptyUserID, errors.Join(m.TgBotRepoErrors.GettingUserIDErr, err)
	}
	if !resultTuple.UserID.IsValid() {
		return m.EmptyUserID, m.TgBotRepoErrors.NoSuchUserErr
	}
	return resultTuple.UserID, nil
}

func (tr *TgBotRepo) RegisterUser(tgID m.TgUserID, userID m.UserID) error {
	if !tgID.IsValid() {
		return m.InvalidTgUserIDErr
	}
	resp, err := tr.conn.Insert(tgIdSpace, []interface{}{tgID, userID})
	if resp.Code != tarantool.OkCode {
		tr.l.Error("error register user in tg_id space", zap.Error(err), zap.Uint32("tarantool response code", resp.Code), zap.Int("tgID", tgID.Int()))
		if err != nil {
			return errors.Join(m.TgBotRepoErrors.RegisterUserErr, err)
		}
		return m.TgBotRepoErrors.RegisterUserErr
	}
	return nil
}
