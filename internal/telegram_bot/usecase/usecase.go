package usecase

import (
	passmanProto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	tgbot "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot"
	"go.uber.org/zap"
)

type TgBotUsecase struct {
	l          *zap.Logger
	r          tgbot.TgBotRepository
	passmanCli passmanProto.PassmanServiceClient
}

func NewTgBotUsecase(l *zap.Logger, r tgbot.TgBotRepository, passmanCli passmanProto.PassmanServiceClient) *TgBotUsecase {
	return &TgBotUsecase{
		l:          l,
		r:          r,
		passmanCli: passmanCli,
	}
}
