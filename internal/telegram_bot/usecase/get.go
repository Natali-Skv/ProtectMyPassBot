package usecase

import (
	"context"
	"errors"
	m "github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	passmanProto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	tgbot "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot"
	"google.golang.org/grpc/status"
	"strings"
)

func (u *TgBotUsecase) GetCommand(req *m.GetCommandReqU) (resp *m.GetCommandRespU, err error) {
	args := strings.Fields(req.ArgsString)
	if len(args) != tgbot.GetCommand.ArgumentCount {
		return nil, m.TgBotUsecaseErrors.WrongArgCountErr
	}

	userID, err := u.r.GetUserID(req.TgID)

	if err != nil {
		switch {
		case errors.Is(err, m.TgBotRepoErrors.NoSuchUserErr):
			return nil, errors.Join(m.TgBotUsecaseErrors.NoSuchUserErr, err)
		default:
			return nil, errors.Join(m.TgBotUsecaseErrors.GetingUserIDUnknownErr, err)
		}
	}

	credentials, err := u.passmanCli.GetCredentials(context.Background(), &passmanProto.GetReq{UserID: userID.Int64(), ServiceName: args[tgbot.GetCommandServiceArgIdx]})
	switch {
	case err == nil:
		return &m.GetCommandRespU{Service: credentials.ServiceName, Login: credentials.Login, Password: credentials.Password}, nil
	case int(status.Code(err)) == m.PassmanHandlerErrors.NoSuchUserOrServiceErr.Code:
		u.l.Debug("37")
		return nil, errors.Join(m.TgBotUsecaseErrors.NoSuchCredsErr)
	default:
		return nil, errors.Join(m.TgBotUsecaseErrors.GettingUserCredsErr, err)
	}
}
