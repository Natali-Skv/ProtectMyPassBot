package repository

import (
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	tgbot "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
)

const (
	tgIdSpace = "tg_id"
	primary   = "primary"
)

type tgIdTuple struct {
	TgID   models.TgUserID
	UserID models.UserID
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

func (tr *TgBotRepo) GetUserID(tgID models.TgUserID) (userID models.UserID, err error) {
	resultTuple := tgIdTuple{}
	err = tr.conn.GetTyped(tgIdSpace, primary, []interface{}{tgID}, &resultTuple)
	if err != nil {
		return models.EmptyUserID, errors.Join(tgbot.GettingUserIDErr, err)
	}
	if !resultTuple.UserID.IsValid() {
		return models.EmptyUserID, tgbot.NoSuchUserErr
	}
	return resultTuple.UserID, nil
}

func (tr *TgBotRepo) RegisterUser(tgID models.TgUserID, userID models.UserID) error {
	if !tgID.IsValid() {
		return models.InvalidTgUserIDErr
	}
	insertTuple := tgIdTuple{TgID: tgID, UserID: userID}
	resp, err := tr.conn.Insert(tgIdSpace, &insertTuple)
	if resp.Code != tarantool.OkCode {
		tr.l.Error("error register user in tg_id space", zap.Error(err), zap.Uint32("tarantool response code", resp.Code), zap.Any("insert tuple", insertTuple))
		if err != nil {
			return errors.Join(tgbot.RegisterUserErr, err)
		}
		return tgbot.RegisterUserErr
	}
	return nil
}
