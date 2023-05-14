package usecase

import "github.com/Natali-Skv/ProtectMyPassBot/internal/models"

func (u *TgBotUsecase) RegisterUser(tgID models.TgUserID) (userID models.UserID, err error) {
	if !tgID.IsValid() {
		return models.EmptyUserID, models.InvalidTgUserIDErr
	}

	//u.passmanCli.

	//TODO
	// register user in passman returning userid

	//u.r.RegisterUser()
	return models.EmptyUserID, nil
}
