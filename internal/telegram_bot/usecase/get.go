package usecase

import (
	"context"
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	passmanProto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	tgbot "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot"
	"google.golang.org/grpc/status"
	"strings"
)

func (u *TgBotUsecase) GetCommand(req *models.GetCommandReqU) (resp *models.GetCommandRespU, err error) {
	args := strings.Fields(req.ArgsString)
	if len(args) != tgbot.GetCommand.ArgumentCount {
		return nil, models.TgBotUsecaseErrors.WrongArgCountErr
	}

	userID, err := u.r.GetUserID(req.TgID)

	if err != nil {
		switch {
		case errors.Is(err, models.TgBotRepoErrors.NoSuchUserErr):
			userID, err = u.RegisterUser(req.TgID)
			if err != nil {
				return nil, errors.Join(models.TgBotUsecaseErrors.NoSuchUserErr, err)
			}
		default:
			return nil, errors.Join(models.TgBotUsecaseErrors.GetingUserIDUnknownErr, err)
		}
	}

	credentials, err := u.passmanCli.GetCredentials(context.Background(), &passmanProto.GetReq{UserID: userID.Int64(), ServiceName: args[tgbot.GetCommandServiceArgumentNumber]})
	switch {
	case err == nil:
		return &models.GetCommandRespU{Service: credentials.ServiceName, Login: credentials.Login, Password: credentials.Password}, nil
	case int(status.Code(err)) == models.PassmanHandlerErrors.NoSuchUserOrServiceErr.Code:
		u.l.Debug("37")
		return nil, errors.Join(models.TgBotUsecaseErrors.NoSuchCredsErr)
	default:
		return nil, errors.Join(models.TgBotUsecaseErrors.GettingUserCredsErr, err)
	}
}

func (u *TgBotUsecase) RegisterUser(tgID models.TgUserID) (userID models.UserID, err error) {
	if !tgID.IsValid() {
		return models.EmptyUserID, models.InvalidTgUserIDErr
	}
	//TODO
	// register user in passman returning userid

	//u.r.RegisterUser()
	return models.EmptyUserID, nil
}
