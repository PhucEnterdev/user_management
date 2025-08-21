package services

import (
	"enterdev.com.vn/user_management/internal/models"
	"enterdev.com.vn/user_management/internal/repository"
	"enterdev.com.vn/user_management/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func (us *UserServiceImpl) GetAllUser() ([]models.User, error) {
	users, err := us.repo.FindAll()
	if err != nil {
		return nil, utils.WrapError(err, "failed to fetch users", utils.ErrCodeInternal)
	}
	return users, nil
}

func (us *UserServiceImpl) CreateUser(user models.User) (models.User, error) {
	user.Email = utils.NormalizeString(user.Email)
	if newUser, existed := us.repo.FindByEmail(user.Email); existed {
		return newUser, utils.NewError("email already exist", utils.ErrCodeConflict)
	}
	user.UUID = uuid.NewString()
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, utils.WrapError(err, "failed to hash password", utils.ErrCodeInternal)
	}
	user.Password = string(hashedPass)
	if err := us.repo.Create(user); err != nil {
		return models.User{}, utils.WrapError(err, "failed to create user", utils.ErrCodeInternal)
	}
	return user, nil
}

func (us *UserServiceImpl) GetUserByUUID(uuid string) (models.User, error) {
	user, found := us.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}
	return user, nil
}

func (us *UserServiceImpl) UpdateUser() {

}

func (us *UserServiceImpl) DeleteUser() {

}
