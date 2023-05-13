package grpc

import (
	"context"
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/passman"
	proto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PassmanHandler struct {
	u passman.PassmanUsecase
	l *zap.Logger
	proto.UnimplementedPassmanServiceServer
}

func NewPassmanHandler(l *zap.Logger, u passman.PassmanUsecase) *PassmanHandler {
	return &PassmanHandler{
		l: l,
		u: u,
	}
}

func (g *PassmanHandler) GetCredentials(ctx context.Context, req *proto.GetReq) (*proto.ServiceCredentials, error) {
	resp, err := g.u.Get(models.GetReqU{
		UserID:  models.UserID(req.UserID),
		Service: req.ServiceName,
	})
	if err != nil {
		switch {
		case errors.Is(err, models.PassmanUsecaseErrors.NoSuchUserOrServiceErr):
			return nil, status.Error(codes.Code(models.PassmanHandlerErrors.NoSuchUserOrServiceErr.Code), errors.Join(models.PassmanHandlerErrors.NoSuchUserOrServiceErr.Error, err).Error())
		default:
			g.l.Error("error getting credentials", zap.Error(err))
			return nil, status.Error(codes.Code(models.PassmanHandlerErrors.UnknownGettingUserCredsErr.Code), errors.Join(models.PassmanHandlerErrors.UnknownGettingUserCredsErr.Error, err).Error())
		}
	}
	return &proto.ServiceCredentials{ServiceName: resp.Service, Login: resp.Login, Password: resp.Password}, nil
}
