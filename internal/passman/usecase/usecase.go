package usecase

import (
	"github.com/Natali-Skv/ProtectMyPassBot/internal/passman"
	"go.uber.org/zap"
)

type PassmanUsecase struct {
	l *zap.Logger
	r passman.PassmanRepository
}

func NewPassmanUsecase(l *zap.Logger, r passman.PassmanRepository) *PassmanUsecase {
	return &PassmanUsecase{l: l, r: r}
}

func (u *PassmanUsecase) Get(req passman.GetReqU) (passman.GetRespU, error) {
	u.l.Debug("req", zap.Any("user_id", req.UserID), zap.String("servicename", req.Service))
	resp, err := u.r.Get(passman.GetReqR{
		UserID:  req.UserID,
		Service: req.Service,
	})
	return passman.GetRespU{
		Service:  resp.Service,
		Login:    resp.Login,
		Password: resp.Password,
	}, err
}
