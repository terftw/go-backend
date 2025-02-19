package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/terftw/go-backend/internal/api/dto"
	"github.com/terftw/go-backend/internal/api/models"
	"github.com/terftw/go-backend/internal/customerrors"
	"github.com/terftw/go-backend/internal/db/repositories"
)

type UserHandler struct {
	UserRepository *repositories.UserRepository
}

func NewUserHandler(ur *repositories.UserRepository) *UserHandler {
	return &UserHandler{
		UserRepository: ur,
	}
}

func (uh *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user := r.Context().Value("user").(*models.User)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		customerrors.ServerErrorResponse(w, r, err)
		return
	}
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	var updates dto.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		customerrors.ErrorResponse(w, r, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := uh.UserRepository.Update(user); err != nil {
		customerrors.ServerErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
