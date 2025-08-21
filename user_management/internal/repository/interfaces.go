package repository

import "enterdev.com.vn/user_management/internal/models"

type UserRepository interface {
	FindAll() ([]models.User, error)
	Create(models.User) error
	FindByUUID(uuid string) (models.User, bool)
	Update()
	Delete()
	FindByEmail(email string) (models.User, bool)
}
