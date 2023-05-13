package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	passmanProto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	tgbot "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot"
	"go.uber.org/zap"
	"strings"
)

func (u *TgBotUsecase) GetCommand(req tgbot.GetCommandReqU) (resp string, err error) {
	args := strings.Fields(req.ArgsString)
	u.l.Debug("args", zap.Strings("args", args))
	if len(args) != models.GetCommand.ArgumentCount {
		return models.GetCommand.Usage, tgbot.WrongArgCountErr
	}

	userID, err := u.r.GetUserID(req.TgID)

	if err != nil {
		switch {
		case errors.Is(err, tgbot.NoSuchUserErr):
			userID, err = u.RegisterUser(req.TgID)
			if err != nil {
				//	TODO
				return "", errors.Join(err)
			}
		default:
			// TODO
			return "", errors.Join(err)
		}
	}

	credentials, err := u.passmanCli.GetCredentials(context.TODO(), &passmanProto.GetReq{UserID: userID.Int64(), ServiceName: args[models.GetCommandServiceArgumentNumber]})
	if err != nil {
		//	TODO
		return "", errors.Join()
	}
	return fmt.Sprintf(models.GetCommand.RespFmtString, credentials.ServiceName, credentials.Login, credentials.Password), nil
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
