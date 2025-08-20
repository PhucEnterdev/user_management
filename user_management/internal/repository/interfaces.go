package repository

import "enterdev.com.vn/user_management/internal/models"

type UserRepository interface {
	FindAll()
	Create()
	FindByUUID()
	Update()
	Delete()
	FindByEmail(email string) (models.User, bool)
}
