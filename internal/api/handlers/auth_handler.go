package handlers

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/terftw/go-backend/internal/api/dto"
	"github.com/terftw/go-backend/internal/api/models"
	"github.com/terftw/go-backend/internal/db/repositories"
	"github.com/terftw/go-backend/internal/utils"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	oauthConfig *oauth2.Config
	userRepo    *repositories.UserRepository
	privateKey  *rsa.PrivateKey
}

func NewAuthHandler(
	oauthConfig *oauth2.Config,
	userRepo *repositories.UserRepository,
	privateKey *rsa.PrivateKey,
) *AuthHandler {
	return &AuthHandler{
		oauthConfig: oauthConfig,
		userRepo:    userRepo,
		privateKey:  privateKey,
	}
}

// InitiateGoogleOAuth starts the OAuth flow
func (h *AuthHandler) InitiateGoogleOAuth(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	state := utils.GenerateRandomState()

	// Store state in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,                // Set to false for development
		SameSite: http.SameSiteLaxMode, // Changed to Lax for development
		Path:     "/",                  // Add path to ensure cookie is sent for all routes
	})

	// Redirect to Google
	url := h.oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "State cookie not found", http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("state") != stateCookie.Value {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	code := r.URL.Query().Get("code")
	token, err := h.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Get user info from Google
	googleUser, err := h.getUserInfo(token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Find or create user in our database
	user, err := h.userRepo.FindOrCreateByGoogleID(googleUser)
	if err != nil {
		http.Error(w, "Failed to process user", http.StatusInternalServerError)
		return
	}

	// Create session token
	sessionToken, err := h.createSession(user)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24 hours
	})

	// Redirect to frontend
	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

// getUserInfo fetches the user's info from Google
func (h *AuthHandler) getUserInfo(accessToken string) (*dto.GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo dto.GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &userInfo, nil
}

// createSession creates a JWT session token
func (h *AuthHandler) createSession(user *models.User) (string, error) {
	// Create JWT token with RS256
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(), // Issued at
		"iss":     "chisel-go",       // Issuer
	})

	// Sign token with private key
	tokenString, err := token.SignedString(h.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// Logout handler
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Redirect to home
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
