package usecase

import (
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	"go.uber.org/zap"
)

type PassmanUsecase struct {
	l *zap.Logger
	r models.PassmanRepository
}

func NewPassmanUsecase(l *zap.Logger, r models.PassmanRepository) *PassmanUsecase {
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

func (u *PassmanUsecase) Register() (models.UserID, error) {
	userID, err := u.r.Register()
	switch {
	case err == nil:
		return userID, nil
	default:
		return models.EmptyUserID, errors.Join(models.PassmanUsecaseErrors.AddNewUserToUserCredsErr, err)
	}
}

func (u *PassmanUsecase) Set(req models.SetReqU) error {
	err := u.r.AddCredentials(models.AddCredsReqR{
		UserID: req.UserID,
		Data:   models.AddCredsData{Service: req.Data.Service, Login: req.Data.Login, Password: req.Data.Password},
	})
	switch {
	case err == nil:
		return nil
	case errors.Is(err, models.PassmanRepoErrors.NoSuchUserErr):
		return errors.Join(models.PassmanUsecaseErrors.NoSuchUserErr, err)
	default:
		return errors.Join(models.PassmanUsecaseErrors.SetUserCredsErr, err)
	}
}

func (u *PassmanUsecase) Del(req models.DeleteCredsReqU) error {
	err := u.r.DeleteCreds(models.DeleteCredsReqR{UserID: req.UserID, Service: req.Service})
	switch {
	case err == nil:
		return nil
	case errors.Is(err, models.PassmanRepoErrors.NoSuchUserErr):
		return errors.Join(models.PassmanUsecaseErrors.NoSuchUserErr, err)
	case errors.Is(err, models.PassmanRepoErrors.NoSuchServiceErr):
		return errors.Join(models.PassmanUsecaseErrors.NoSuchServiceErr, err)
	default:
		return errors.Join(models.PassmanUsecaseErrors.DeleteUserCredsErr, err)
	}
}
