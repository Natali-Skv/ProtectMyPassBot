package passman

import "github.com/Natali-Skv/ProtectMyPassBot/internal/models"

type PassmanUsecase interface {
	Get(req models.GetReqU) (models.GetRespU, error)
}
