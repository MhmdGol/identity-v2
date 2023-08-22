package controller

import (
	"context"
	"identity-v2/internal/model"
	userapiv1 "identity-v2/internal/proto/userapi/v1"
	"identity-v2/internal/service"

	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserController struct {
	userSvc service.UserService
	userapiv1.UnimplementedUserServiceServer
}

var _ userapiv1.UserServiceServer = (*UserController)(nil)

func NewUserController(
	userSvc service.UserService,
) *UserController {
	return &UserController{
		userSvc: userSvc,
	}
}

func (uc *UserController) CreateUser(ctx context.Context, req *userapiv1.CreateUserRequest) (*userapiv1.CreateUserResponse, error) {
	err := uc.userSvc.Create(ctx, model.RawUser{
		UUN:      req.Uun,
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
		Status:   req.Status,
	})
	if err != nil {
		return &userapiv1.CreateUserResponse{}, status.Error(codes.Code(code.Code_UNKNOWN), "not created")
	}

	return &userapiv1.CreateUserResponse{}, status.Error(codes.Code(code.Code_OK), "created")
}
