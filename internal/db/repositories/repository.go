package repositories

import "gorm.io/gorm"

type Repositories struct {
	UserRepository *UserRepository
}

func NewRepository(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepository: NewUserRepository(db),
	}
}
