package service

import (
	"identity-v2/internal/repository"
	"identity-v2/internal/service"

	"github.com/uptrace/bun"
)

type LoginAttemptService struct {
	loginAttemptRepo repository.LoginAttemptRepo
}

var _ service.LoginAttemptService = (*LoginAttemptService)(nil)

func NewLoginAttempService(

) *LoginAttemptService {
	return &LoginAttemptService{

	}
}

func (las *LoginAttemptService) 