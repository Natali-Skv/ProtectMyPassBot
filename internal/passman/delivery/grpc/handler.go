package grpc

import (
	"context"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/passman"
	proto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
)

type grpcPassmanHandler struct {
	u passman.PassmanUsecase
	proto.UnimplementedPassmanServiceServer
}

func NewPassmanHandler(u passman.PassmanUsecase) *grpcPassmanHandler {
	return &grpcPassmanHandler{
		u: u,
	}
}

func (g *grpcPassmanHandler) GetCredentials(ctx context.Context, req *proto.GetReq) (*proto.ServiceCredentials, error) {
	return &proto.ServiceCredentials{ServiceName: "VK", Credentials: []*proto.Credentials{&proto.Credentials{Login: "natali", Password: "maypass"}}}, nil
}
