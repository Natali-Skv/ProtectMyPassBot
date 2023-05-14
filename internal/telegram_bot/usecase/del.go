package usecase

import (
	"context"
	"errors"
	m "github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	proto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	tgbot "github.com/Natali-Skv/ProtectMyPassBot/internal/telegram_bot"
	"google.golang.org/grpc/status"
	"strings"
)

func (u *TgBotUsecase) DelCommand(req *m.DelCommandReqU) (serviceName string, err error) {
	args := strings.Fields(req.ArgsString)
	if len(args) != tgbot.DelCommand.ArgumentCount {
		return "", m.TgBotUsecaseErrors.WrongArgCountErr
	}

	serviceName = args[tgbot.DelCommandServiceArgIdx]

	userID, err := u.r.GetUserID(req.TgID)

	switch {
	case err == nil:
	case errors.Is(err, m.TgBotRepoErrors.NoSuchUserErr):
		return "", errors.Join(m.TgBotUsecaseErrors.NoSuchUserErr, err)
	default:
		return "", errors.Join(m.TgBotUsecaseErrors.GetingUserIDUnknownErr, err)
	}

	_, err = u.passmanCli.DelCredentials(context.Background(), &proto.DelReq{UserID: userID.Int64(), ServiceName: serviceName})

	switch {
	case err == nil:
		return serviceName, nil
	case int(status.Code(err)) == m.PassmanHandlerErrors.NoSuchServiceErr.Code:
		return "", errors.Join(m.TgBotUsecaseErrors.NoSuchCredsErr, err)
	default:
		return "", errors.Join(m.TgBotUsecaseErrors.DelUserCredsErr, err)
	}
}
