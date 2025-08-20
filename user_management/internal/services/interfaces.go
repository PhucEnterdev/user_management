package services

import "enterdev.com.vn/user_management/internal/models"

type UserService interface {
	GetAllUser()
	CreateUser(user models.User)
	GetUserByUUID()
	UpdateUser()
	DeleteUser()
}
