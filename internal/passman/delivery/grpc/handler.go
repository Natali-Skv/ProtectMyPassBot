package grpc

import (
	"context"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/passman"
	proto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	"go.uber.org/zap"
)

type grpcPassmanHandler struct {
	u passman.PassmanUsecase
	l *zap.Logger
	proto.UnimplementedPassmanServiceServer
}

func NewPassmanHandler(l *zap.Logger, u passman.PassmanUsecase) *grpcPassmanHandler {
	return &grpcPassmanHandler{
		l: l,
		u: u,
	}
}

func (g *grpcPassmanHandler) GetCredentials(ctx context.Context, req *proto.GetReq) (*proto.ServiceCredentials, error) {
	// go to usecase
	g.l.Debug("req", zap.Any("user_id", req.UserID), zap.String("servicename", req.ServiceName))
	resp, err := g.u.Get(passman.GetReqU{
		UserID:  models.UserID(req.UserID),
		Service: req.ServiceName,
	})
	return &proto.ServiceCredentials{ServiceName: resp.Service, Login: resp.Login, Password: resp.Password}, err
}
