package handlers

import (
	"crypto/rsa"

	"github.com/terftw/go-backend/internal/db/repositories"
	"golang.org/x/oauth2"
)

type Handlers struct {
	User *UserHandler
	Auth *AuthHandler
}

func NewHandlers(
	userRepo *repositories.UserRepository,
	oauthConfig *oauth2.Config,
	privateKey *rsa.PrivateKey,
) *Handlers {
	return &Handlers{
		User: NewUserHandler(userRepo),
		Auth: NewAuthHandler(oauthConfig, userRepo, privateKey),
	}
}
