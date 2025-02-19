package repositories

import (
	"github.com/terftw/go-backend/internal/api/dto"
	"github.com/terftw/go-backend/internal/api/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindOrCreateByGoogleID will either find existing user or create new one
func (r *UserRepository) FindOrCreateByGoogleID(
	googleUserInfo *dto.GoogleUserInfo,
) (*models.User, error) {
	var user models.User

	// Try to find existing user by GoogleID
	err := r.db.Where("google_id = ?", googleUserInfo.ID).First(&user).Error
	if err == nil {
		// User found, update their info
		user.Name = googleUserInfo.Name
		user.Picture = googleUserInfo.Picture

		if err := r.db.Save(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	// User not found, create new one
	newUser := &models.User{
		Email:    googleUserInfo.Email,
		Name:     googleUserInfo.Name,
		GoogleID: googleUserInfo.ID,
		Picture:  googleUserInfo.Picture,
	}

	if err := r.db.Create(newUser).Error; err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}
