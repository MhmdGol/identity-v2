package controller

import (
	"context"
	"identity-v2/internal/model"
	"identity-v2/pkg/jwt"

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
