package controller

import (
	"context"
	"identity-v2/internal/model"
	"identity-v2/internal/service"
	"identity-v2/pkg/jwt"

	"github.com/casbin/casbin/v2"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ExtractHeader(ctx context.Context, j *jwt.JwtToken) (model.TokenClaim, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return model.TokenClaim{}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	token, found := md["authorization"]
	if !found || len(token) != 1 {
		return model.TokenClaim{}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	tc, err := j.ExtractClaims(model.JwtToken(token[0]))
	if err != nil {
		return model.TokenClaim{}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	return tc, nil
}

func LoginPermissionCheck(ctx context.Context, j *jwt.JwtToken, authSvc service.AuthService, e *casbin.Enforcer, self bool, obj, act string) (int64, string, error) {
	tc, err := ExtractHeader(ctx, j)
	if err != nil {
		return 0, "", err
	}

	loggedin, err := authSvc.CheckSession(ctx, tc.ID)
	if err != nil || !loggedin {
		return 0, "", status.Error(codes.Code(code.Code_INTERNAL), "login first")
	}

	if !self {
		ok, err := e.Enforce(tc.Email, obj, act)
		if err != nil || !ok {
			return 0, "", status.Error(codes.Code(code.Code_PERMISSION_DENIED), "not allowed")
		}
	}

	return int64(tc.ID), tc.Email, nil
}
