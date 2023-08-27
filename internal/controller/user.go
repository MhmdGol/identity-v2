package controller

import (
	"context"
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
	tc, err := ExtractHeader(ctx, uc.jwt)
	if err != nil {
		return &userapiv1.CreateUserResponse{}, err
	}

	loggedin, err := uc.authSvc.CheckSession(ctx, tc.ID)
	if err != nil || !loggedin {
		return &userapiv1.CreateUserResponse{}, status.Error(codes.Code(code.Code_INTERNAL), "login first")
	}

	ok, err := uc.e.Enforce(tc.Email, "users", "create")
	if err != nil || !ok {
		return &userapiv1.CreateUserResponse{}, status.Error(codes.Code(code.Code_PERMISSION_DENIED), "not allowed")
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
