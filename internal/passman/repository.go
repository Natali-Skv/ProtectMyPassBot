package passman

import (
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
)

type PassmanRepository interface {
	Get(req models.GetReqR) (models.GetRespR, error)
}
