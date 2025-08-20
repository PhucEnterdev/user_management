package services

import (
	"log"

	"enterdev.com.vn/user_management/internal/models"
	"enterdev.com.vn/user_management/internal/repository"
	"enterdev.com.vn/user_management/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserServiceImpl struct {
	// vì UserRepo trả về một *InMemoryUserRepository
	// nên ta thêm *InMemoryUserRepository vào struct UserServiceImpl
	repo repository.UserRepository
}

func NewUserServiceImpl(repo repository.UserRepository) UserService {
	return &UserServiceImpl{
		repo: repo,
	}
}

func (us *UserServiceImpl) GetAllUser() {
	us.repo.FindAll()
	log.Println("GetAllUser into UserServiceImpl")
}

func (us *UserServiceImpl) CreateUser(user models.User) models.User {
	user.Email = utils.NormalizeString(user.Email)
	if newUser, existed := us.repo.FindByEmail(user.Email); existed {
		return newUser, gin.H{""}
	}
}

func (us *UserServiceImpl) GetUserByUUID() {

}

func (us *UserServiceImpl) UpdateUser() {

}

func (us *UserServiceImpl) DeleteUser() {

}
