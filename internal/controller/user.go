package controller

import (
	"context"
	"fmt"
	"identity-v2/internal/model"
	userapiv1 "identity-v2/internal/proto/userapi/v1"
	"identity-v2/internal/service"
	"identity-v2/pkg/jwt"

	"github.com/casbin/casbin/v2"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserController struct {
	userSvc service.UserService
	authSvc service.AuthService
	jwt     *jwt.JwtToken
	e       *casbin.Enforcer
	userapiv1.UnimplementedUserServiceServer
}

var _ userapiv1.UserServiceServer = (*UserController)(nil)

func NewUserController(
	userSvc service.UserService,
	authSvc service.AuthService,
	jwt *jwt.JwtToken,
	e *casbin.Enforcer,
) *UserController {
	return &UserController{
		userSvc: userSvc,
		authSvc: authSvc,
		jwt:     jwt,
		e:       e,
	}
}

func (uc *UserController) CreateUser(ctx context.Context, req *userapiv1.CreateUserRequest) (*userapiv1.CreateUserResponse, error) {
	// logged in and active seesion
	_, _, err := LoginPermissionCheck(ctx, uc.jwt, uc.authSvc, uc.e, false, "users", "create")
	if err != nil {
		return &userapiv1.CreateUserResponse{}, err
	}

	err = uc.userSvc.Create(ctx, model.RawUser{
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

func (uc *UserController) SetTOTP(ctx context.Context, req *userapiv1.SetTOTPRequest) (*userapiv1.SetTOTPResponse, error) {
	_, email, err := LoginPermissionCheck(ctx, uc.jwt, uc.authSvc, uc.e, true, "", "")
	if err != nil {
		return &userapiv1.SetTOTPResponse{}, err
	}

	secret, err := uc.userSvc.SetTOTP(ctx, email)
	if err != nil {
		fmt.Println(err)
		return &userapiv1.SetTOTPResponse{}, status.Error(codes.Code(code.Code_INTERNAL), "something went wrong")
	}

	return &userapiv1.SetTOTPResponse{
		TotpSecret: secret,
	}, status.Error(codes.Code(code.Code_OK), "all good")
}
