package usecase

import (
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
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

func (u *PassmanUsecase) Get(req models.GetReqU) (models.GetRespU, error) {
	u.l.Debug("req", zap.Any("user_id", req.UserID), zap.String("servicename", req.Service))
	resp, err := u.r.Get(models.GetReqR{
		UserID:  req.UserID,
		Service: req.Service,
	})

	if err != nil {
		switch {
		case errors.Is(err, models.PassmanRepoErrors.NoSuchUserOrServiceInDBErr):
			return models.GetRespU{}, errors.Join(models.PassmanUsecaseErrors.NoSuchUserOrServiceErr)
		default:
			return models.GetRespU{}, errors.Join(models.PassmanUsecaseErrors.UnknownGettingUserCredsErr, err)
		}
	}

	return models.GetRespU{
		Service:  resp.Service,
		Login:    resp.Login,
		Password: resp.Password,
	}, nil
}
