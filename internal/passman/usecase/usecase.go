package usecase

import (
	"github.com/Natali-Skv/ProtectMyPassBot/internal/passman"
)

type PassmanUsecase struct {
}

func NewPassmanUsecase() *PassmanUsecase {
	return &PassmanUsecase{}
}

func (u *PassmanUsecase) Get(req passman.GetReqU) (passman.GetRespU, error) {
	return passman.GetRespU{}, nil
}
