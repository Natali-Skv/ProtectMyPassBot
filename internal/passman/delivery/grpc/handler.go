package grpc

import (
	"context"
	"errors"
	"github.com/Natali-Skv/ProtectMyPassBot/internal/models"
	proto "github.com/Natali-Skv/ProtectMyPassBot/internal/passman/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PassmanHandler struct {
	u models.PassmanUsecase
	l *zap.Logger
	proto.UnimplementedPassmanServiceServer
}

func NewPassmanHandler(l *zap.Logger, u models.PassmanUsecase) *PassmanHandler {
	return &PassmanHandler{
		l: l,
		u: u,
	}
}

func (h *PassmanHandler) GetCredentials(ctx context.Context, req *proto.GetReq) (*proto.ServiceCredentials, error) {
	resp, err := h.u.Get(models.GetReqU{
		UserID:  models.UserID(req.UserID),
		Service: req.ServiceName,
	})
	switch {
	case err == nil:
		return &proto.ServiceCredentials{ServiceName: resp.Service, Login: resp.Login, Password: resp.Password}, nil
	case errors.Is(err, models.PassmanUsecaseErrors.NoSuchUserOrServiceErr):
		return nil, status.Error(codes.Code(models.PassmanHandlerErrors.NoSuchUserOrServiceErr.Code), errors.Join(models.PassmanHandlerErrors.NoSuchUserOrServiceErr.Error, err).Error())
	default:
		h.l.Error("error getting credentials", zap.Error(err))
		return nil, status.Error(codes.Code(models.PassmanHandlerErrors.UnknownGettingUserCredsErr.Code), errors.Join(models.PassmanHandlerErrors.UnknownGettingUserCredsErr.Error, err).Error())
	}
}

func (h *PassmanHandler) Register(ctx context.Context) (*proto.RegisterResp, error) {
	userID, err := h.u.Register()
	switch {
	case err == nil:
		return &proto.RegisterResp{
			UserID: userID.Int64(),
		}, nil
	default:
		return nil, status.Error(codes.Code(models.PassmanHandlerErrors.RegisterUserErr.Code), errors.Join(models.PassmanHandlerErrors.RegisterUserErr.Error, err).Error())
	}
}

func (h *PassmanHandler) SetCredentials(ctx context.Context, req *proto.SetReq) error {
	err := h.u.Set(models.SetReqU{
		UserID: models.UserID(req.UserID),
		Data:   models.AddCredsData{Login: req.Data.Login, Password: req.Data.Password, Service: req.Data.ServiceName},
	})
	switch {
	case err == nil:
		return nil
	case errors.Is(err, models.PassmanUsecaseErrors.NoSuchUserErr):
		return status.Error(codes.Code(models.PassmanHandlerErrors.NoSuchUserErr.Code), errors.Join(models.PassmanHandlerErrors.NoSuchUserErr.Error, err).Error())
	default:
		return status.Error(codes.Code(models.PassmanHandlerErrors.SetUserCredsErr.Code), errors.Join(models.PassmanHandlerErrors.SetUserCredsErr.Error, err).Error())
	}
}

func (h *PassmanHandler) DelCredentials(ctx context.Context, req *proto.DelReq) error {
	err := h.u.Del(models.DeleteCredsReqU{UserID: models.UserID(req.UserID), Service: req.ServiceName})
	switch {
	case err == nil:
		return nil
	case errors.Is(err, models.PassmanUsecaseErrors.NoSuchUserErr):
		return status.Error(codes.Code(models.PassmanHandlerErrors.NoSuchUserErr.Code), errors.Join(models.PassmanHandlerErrors.NoSuchUserErr.Error, err).Error())
	case errors.Is(err, models.PassmanUsecaseErrors.NoSuchServiceErr):
		return status.Error(codes.Code(models.PassmanHandlerErrors.NoSuchServiceErr.Code), errors.Join(models.PassmanHandlerErrors.NoSuchServiceErr.Error, err).Error())
	default:
		return status.Error(codes.Code(models.PassmanHandlerErrors.DelUserCredsErr.Code), errors.Join(models.PassmanHandlerErrors.DelUserCredsErr.Error, err).Error())
	}
}
