package grpc

import (
	"context"
	"errors"
	m "github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	proto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PassmanHandler struct {
	u m.PassmanUsecase
	l *zap.Logger
	proto.UnimplementedPassmanServiceServer
}

func NewPassmanHandler(l *zap.Logger, u m.PassmanUsecase) *PassmanHandler {
	return &PassmanHandler{
		l: l,
		u: u,
	}
}

func (h *PassmanHandler) GetCredentials(ctx context.Context, req *proto.GetReq) (*proto.ServiceCredentials, error) {
	resp, err := h.u.Get(m.GetReqU{
		UserID:  m.UserID(req.UserID),
		Service: req.ServiceName,
	})
	switch {
	case err == nil:
		return &proto.ServiceCredentials{ServiceName: resp.Service, Login: resp.Login, Password: resp.Password}, nil
	case errors.Is(err, m.PassmanUsecaseErrors.NoSuchUserOrServiceErr):
		return nil, status.Error(codes.Code(m.PassmanHandlerErrors.NoSuchUserOrServiceErr.Code), errors.Join(m.PassmanHandlerErrors.NoSuchUserOrServiceErr.Error, err).Error())
	default:
		h.l.Error("error getting credentials", zap.Error(err))
		return nil, status.Error(codes.Code(m.PassmanHandlerErrors.UnknownGettingUserCredsErr.Code), errors.Join(m.PassmanHandlerErrors.UnknownGettingUserCredsErr.Error, err).Error())
	}
}

func (h *PassmanHandler) RegisterUser(ctx context.Context, _ *empty.Empty) (*proto.RegisterResp, error) {
	userID, err := h.u.Register()
	switch {
	case err == nil:
		return &proto.RegisterResp{
			UserID: userID.Int64(),
		}, nil
	default:
		return nil, status.Error(codes.Code(m.PassmanHandlerErrors.RegisterUserErr.Code), errors.Join(m.PassmanHandlerErrors.RegisterUserErr.Error, err).Error())
	}
}

func (h *PassmanHandler) SetCredentials(ctx context.Context, req *proto.SetReq) (*empty.Empty, error) {
	err := h.u.Set(m.SetReqU{
		UserID: m.UserID(req.UserID),
		Data:   m.AddCredsData{Login: req.Data.Login, Password: req.Data.Password, Service: req.Data.ServiceName},
	})
	switch {
	case err == nil:
		return &empty.Empty{}, nil
	case errors.Is(err, m.PassmanUsecaseErrors.NoSuchUserErr):
		return &empty.Empty{}, status.Error(codes.Code(m.PassmanHandlerErrors.NoSuchUserErr.Code), errors.Join(m.PassmanHandlerErrors.NoSuchUserErr.Error, err).Error())
	default:
		return &empty.Empty{}, status.Error(codes.Code(m.PassmanHandlerErrors.SetUserCredsErr.Code), errors.Join(m.PassmanHandlerErrors.SetUserCredsErr.Error, err).Error())
	}
}

func (h *PassmanHandler) DelCredentials(ctx context.Context, req *proto.DelReq) (*empty.Empty, error) {
	err := h.u.Del(m.DeleteCredsReqU{UserID: m.UserID(req.UserID), Service: req.ServiceName})
	switch {
	case err == nil:
		return &empty.Empty{}, nil
	case errors.Is(err, m.PassmanUsecaseErrors.NoSuchUserErr):
		return &empty.Empty{}, status.Error(codes.Code(m.PassmanHandlerErrors.NoSuchUserErr.Code), errors.Join(m.PassmanHandlerErrors.NoSuchUserErr.Error, err).Error())
	case errors.Is(err, m.PassmanUsecaseErrors.NoSuchServiceErr):
		return &empty.Empty{}, status.Error(codes.Code(m.PassmanHandlerErrors.NoSuchServiceErr.Code), errors.Join(m.PassmanHandlerErrors.NoSuchServiceErr.Error, err).Error())
	default:
		return &empty.Empty{}, status.Error(codes.Code(m.PassmanHandlerErrors.DelUserCredsErr.Code), errors.Join(m.PassmanHandlerErrors.DelUserCredsErr.Error, err).Error())
	}
}
