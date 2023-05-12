package repository

import (
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
)

const (
	tgIdSpace = "tg_id"
	primary   = "primary"
)

type TgBotRepo struct {
	l    *zap.Logger // TODO он точно нужен? или все логируется выше
	conn *tarantool.Connection
}

func NewTgBotRepo(l *zap.Logger, conn *tarantool.Connection) *TgBotRepo {
	if l == nil || conn == nil {
		return nil
	}
	return &TgBotRepo{l: l, conn: conn}
}

func (tr *TgBotRepo) GetUserID(tgID int) (userID uint64, err error) {
	queryResuln := struct {
		TgID   int
		UserID uint64
	}{}
	err = tr.conn.GetTyped(tgIdSpace, primary, []interface{}{tgID}, &queryResuln)
	tr.l.Debug("intID", zap.Any("intID", queryResuln))
	if err != nil {
		return 0, errors.Join(telegram_bot.GettingUserIDErr, err)
	}
	return queryResuln.UserID, nil
}
