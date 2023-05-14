package usecase

import (
	"errors"
	m "github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	"go.uber.org/zap"
)

type PassmanUsecase struct {
	l *zap.Logger
	r m.PassmanRepository
}

func NewPassmanUsecase(l *zap.Logger, r m.PassmanRepository) *PassmanUsecase {
	return &PassmanUsecase{l: l, r: r}
}

func (u *PassmanUsecase) Get(req m.GetReqU) (m.GetRespU, error) {
	u.l.Debug("req", zap.Any("user_id", req.UserID), zap.String("servicename", req.Service))
	resp, err := u.r.Get(m.GetReqR{
		UserID:  req.UserID,
		Service: req.Service,
	})

	if err != nil {
		switch {
		case errors.Is(err, m.PassmanRepoErrors.NoSuchUserOrServiceInDBErr):
			return m.GetRespU{}, errors.Join(m.PassmanUsecaseErrors.NoSuchUserOrServiceErr)
		default:
			return m.GetRespU{}, errors.Join(m.PassmanUsecaseErrors.UnknownGettingUserCredsErr, err)
		}
	}

	return m.GetRespU{
		Service:  resp.Service,
		Login:    resp.Login,
		Password: resp.Password,
	}, nil
}

func (u *PassmanUsecase) Register() (m.UserID, error) {
	userID, err := u.r.Register()
	switch {
	case err == nil:
		return userID, nil
	default:
		return m.EmptyUserID, errors.Join(m.PassmanUsecaseErrors.AddNewUserToUserCredsErr, err)
	}
}

func (u *PassmanUsecase) Set(req m.SetReqU) error {
	err := u.r.AddCredentials(m.AddCredsReqR{
		UserID: req.UserID,
		Data:   m.AddCredsData{Service: req.Data.Service, Login: req.Data.Login, Password: req.Data.Password},
	})
	switch {
	case err == nil:
		return nil
	case errors.Is(err, m.PassmanRepoErrors.NoSuchUserErr):
		return errors.Join(m.PassmanUsecaseErrors.NoSuchUserErr, err)
	default:
		return errors.Join(m.PassmanUsecaseErrors.SetUserCredsErr, err)
	}
}

func (u *PassmanUsecase) Del(req m.DeleteCredsReqU) error {
	err := u.r.DeleteCreds(m.DeleteCredsReqR{UserID: req.UserID, Service: req.Service})
	switch {
	case err == nil:
		return nil
	case errors.Is(err, m.PassmanRepoErrors.NoSuchUserErr):
		return errors.Join(m.PassmanUsecaseErrors.NoSuchUserErr, err)
	case errors.Is(err, m.PassmanRepoErrors.NoSuchServiceErr):
		return errors.Join(m.PassmanUsecaseErrors.NoSuchServiceErr, err)
	default:
		return errors.Join(m.PassmanUsecaseErrors.DeleteUserCredsErr, err)
	}
}
