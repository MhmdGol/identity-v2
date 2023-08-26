package controller

import (
	"context"
	"fmt"
	"identity-v2/internal/model"
	authapiv1 "identity-v2/internal/proto/authapi/v1"
	"identity-v2/internal/service"
	"identity-v2/pkg/jwt"

	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthController struct {
	authSvc service.AuthService
	jwt     *jwt.JwtToken
	authapiv1.UnimplementedAuthServiceServer
}

var _ authapiv1.AuthServiceServer = (*AuthController)(nil)

func NewAuthController(
	authSvc service.AuthService,
	jwt *jwt.JwtToken,
) *AuthController {
	return &AuthController{
		authSvc: authSvc,
		jwt:     jwt,
	}
}

func (ac *AuthController) Login(ctx context.Context, req *authapiv1.LoginRequest) (*authapiv1.LoginResponse, error) {
	t, err := ac.authSvc.Login(ctx, model.LoginInfo{
		Email:    req.Email,
		Password: req.Password,
		TOTPCode: req.TotpCode,
	})
	fmt.Println(err)
	if err != nil {
		return &authapiv1.LoginResponse{}, status.Error(codes.Code(code.Code_INVALID_ARGUMENT), "login failed")
	}

	return &authapiv1.LoginResponse{
		Token: string(t),
	}, status.Error(codes.Code(code.Code_OK), "logged in")
}

func (ac *AuthController) Logout(ctx context.Context, req *authapiv1.LogoutRequest) (*authapiv1.LogoutResponse, error) {
	tc, err := ExtractHeader(ctx, ac.jwt)
	if err != nil {
		return &authapiv1.LogoutResponse{}, err
	}

	err = ac.authSvc.Logout(ctx, tc.ID)
	if err != nil {
		return &authapiv1.LogoutResponse{}, status.Error(codes.Code(code.Code_UNKNOWN), "logout failed")
	}

	return &authapiv1.LogoutResponse{}, status.Error(codes.Code(code.Code_OK), "logged out")
}
