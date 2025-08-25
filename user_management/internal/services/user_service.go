package services

import (
	"log"
	"strings"

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

func (us *UserServiceImpl) GetAllUser(search string, page int, limit int) ([]models.User, error) {
	users, err := us.repo.FindAll()
	if err != nil {
		return nil, utils.WrapError(err, "failed to fetch users", utils.ErrCodeInternal)
	}
	var filteredUsers []models.User
	log.Println(search)
	if search == "" {
		filteredUsers = users
	} else {
		search = strings.ToLower(search)
		for _, user := range users {
			name := strings.ToLower(user.Name)
			email := strings.ToLower(user.Email)
			if strings.Contains(name, search) || strings.Contains(email, search) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	}

	start := (page - 1) * limit
	if start >= len(filteredUsers) {
		return []models.User{}, nil
	}
	end := start + limit
	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}
	return filteredUsers[start:end], nil
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

func (us *UserServiceImpl) UpdateUser(uuid string, user models.User) (models.User, error) {
	user.Email = utils.NormalizeString(user.Email)
	if u, exist := us.repo.FindByEmail(user.Email); exist && u.UUID != uuid {
		return models.User{}, utils.NewError("email already exist", utils.ErrCodeConflict)
	}
	currentUser, found := us.repo.FindByUUID(uuid)
	if !found {
		return models.User{}, utils.NewError("user not found", utils.ErrCodeNotFound)
	}
	currentUser.Name = user.Name
	currentUser.Age = user.Age
	if user.Email != "" {
		currentUser.Email = user.Email
	}
	currentUser.Status = user.Status
	currentUser.Level = user.Level
	if user.Password != "" {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, utils.WrapError(err, "failed to hash password", utils.ErrCodeInternal)
		}
		currentUser.Password = string(hashedPass)
	}

	// update data
	if err := us.repo.Update(uuid, currentUser); err != nil {
		return models.User{}, utils.WrapError(err, "failed to update user", utils.ErrCodeInternal)
	}
	return currentUser, nil
}

func (us *UserServiceImpl) DeleteUser(uuid string) error {
	if err := us.repo.Delete(uuid); err != nil {
		return utils.NewError("user not found", utils.ErrCodeNotFound)
	}
	return nil
}
