package services

import (
	"enterdev.com.vn/user_management/internal/models"
)

type UserService interface {
	GetAllUser() ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetUserByUUID(uuid string) (models.User, error)
	UpdateUser()
	DeleteUser()
}
