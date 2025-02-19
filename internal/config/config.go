package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/terftw/go-backend/internal/api/handlers"
	"github.com/terftw/go-backend/internal/db/repositories"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthConfig struct {
	GoogleOAuth *oauth2.Config
}
type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	OAuth      OAuthConfig
	PrivateKey *rsa.PrivateKey
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func Load() (*Config, error) {
	privateKey, err := loadRSAKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to load RSA keys: %w", err)
	}

	oAuthConfig, err := loadOAuthConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load OAuth config: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     5432,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		OAuth:      *oAuthConfig,
		PrivateKey: privateKey,
	}, nil
}

func loadRSAKeys() (*rsa.PrivateKey, error) {
	privateKeyPEM, err := base64.StdEncoding.DecodeString(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	// Parse PEM block
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	// Parse private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}

func loadOAuthConfig() (*OAuthConfig, error) {
	return &OAuthConfig{
		GoogleOAuth: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}, nil
}

func (c *Config) InitializeHandlers(r *repositories.Repositories) *handlers.Handlers {
	return handlers.NewHandlers(
		r.UserRepository,
		c.OAuth.GoogleOAuth,
		c.PrivateKey,
	)
}
