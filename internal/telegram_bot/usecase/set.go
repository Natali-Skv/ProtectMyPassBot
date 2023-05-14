package usecase

import (
	"context"
	"errors"
	m "github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	proto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	tgbot "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot"
	"github.com/golang/protobuf/ptypes/empty"
	"strings"
)

func (u *TgBotUsecase) SetCommand(req *m.SetCommandReqU) (serviceName string, err error) {
	args := strings.Fields(req.ArgsString)
	if len(args) != tgbot.SetCommand.ArgumentCount {
		return "", m.TgBotUsecaseErrors.WrongArgCountErr
	}

	userID, err := u.r.GetUserID(req.TgID)

	switch {
	case err == nil:
	case errors.Is(err, m.TgBotRepoErrors.NoSuchUserErr):
		registerResp, err := u.passmanCli.RegisterUser(context.Background(), &empty.Empty{})
		if err != nil {
			return "", errors.Join(m.TgBotUsecaseErrors.RegisterUserErr, err)
		}
		userID = m.UserID(registerResp.UserID)
		err = u.r.RegisterUser(req.TgID, userID)
		if err != nil {
			return "", errors.Join(m.TgBotUsecaseErrors.RegisterUserErr, err)
		}
	default:
		return "", errors.Join(m.TgBotUsecaseErrors.GetingUserIDUnknownErr, err)
	}
	serviceName = args[tgbot.SetCommandServiceArgIdx]
	login := args[tgbot.SetCommandLoginArgIdx]
	password := args[tgbot.SetCommandPasswordArgIdx]
	_, err = u.passmanCli.SetCredentials(context.Background(), &proto.SetReq{UserID: userID.Int64(), Data: &proto.ServiceCredentials{ServiceName: serviceName, Login: login, Password: password}})

	if err != nil {
		return "", errors.Join(m.TgBotUsecaseErrors.SetUserCredsErr, err)
	}
	return serviceName, nil
}
