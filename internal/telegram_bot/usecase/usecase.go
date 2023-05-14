package usecase

import (
	m "github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	passmanProto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	"go.uber.org/zap"
)

type TgBotUsecase struct {
	l          *zap.Logger
	r          m.TgBotRepository
	passmanCli passmanProto.PassmanServiceClient
}

func NewTgBotUsecase(l *zap.Logger, r m.TgBotRepository, passmanCli passmanProto.PassmanServiceClient) *TgBotUsecase {
	return &TgBotUsecase{
		l:          l,
		r:          r,
		passmanCli: passmanCli,
	}
}
