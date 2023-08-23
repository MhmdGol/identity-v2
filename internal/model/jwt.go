package model

type JwtToken string

type TokenClaim struct {
	ID    ID
	Email string
}
