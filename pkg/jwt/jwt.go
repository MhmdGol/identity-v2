package jwt

import (
	"crypto/rsa"
	"fmt"
	"identity-v2/cmd/config"
	"identity-v2/internal/model"
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/dgrijalva/jwt-go"
)

type JwtToken struct {
	SecretKey *rsa.PrivateKey
	PublicKey *rsa.PublicKey
}

func NewJwtHandler(conf config.RSAPair) (*JwtToken, error) {
	privateKeyBytes, err := os.ReadFile(conf.SecretKeyPath)
	if err != nil {
		return &JwtToken{}, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return &JwtToken{}, err
	}

	publicKeyBytes, err := os.ReadFile(conf.PublicKeyPath)
	if err != nil {
		return &JwtToken{}, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return &JwtToken{}, err
	}

	return &JwtToken{
		SecretKey: privateKey,
		PublicKey: publicKey,
	}, nil
}

func (j *JwtToken) MakeToken(c model.TokenClaim) (model.JwtToken, error) {
	claims := jwt.MapClaims{
		"id":    strconv.FormatInt(int64(c.ID), 10),
		"email": c.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(j.SecretKey)

	return model.JwtToken(tokenString), err
}

func (j *JwtToken) ExtractClaims(t model.JwtToken) (model.TokenClaim, error) {
	token, err := jwt.Parse(string(t), func(token *jwt.Token) (interface{}, error) {
		return j.PublicKey, nil
	})
	if err != nil {
		return model.TokenClaim{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		id, ok := claims["id"].(string)
		if !ok {
			return model.TokenClaim{}, fmt.Errorf("invalid token: id not found")
		}
		sfId, _ := snowflake.ParseString(id)

		email, ok := claims["email"].(string)
		if !ok {
			fmt.Println(4)
			return model.TokenClaim{}, fmt.Errorf("invalid token: username not found")
		}
		return model.TokenClaim{
			ID:    model.ID(sfId.Int64()),
			Email: email,
		}, nil
	}
	return model.TokenClaim{}, fmt.Errorf("invalid token")
}
