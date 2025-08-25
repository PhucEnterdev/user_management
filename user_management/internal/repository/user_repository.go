package repository

import (
	"fmt"
	"slices"

	"enterdev.com.vn/user_management/internal/models"
)

// Bài này sử dụng slice để lưu data
// Trong TH sử dụng DB thì ta có thể thêm SQLUserRepository hoặc PostgreSQL...
type InMemoryUserRepositoryImpl struct {
	users []models.User
}

func NewInMemoryUserRepositoryImpe() UserRepository {
	return &InMemoryUserRepositoryImpl{
		users: make([]models.User, 0),
	}
}

func (ur *InMemoryUserRepositoryImpl) FindAll() ([]models.User, error) {
	return ur.users, nil
}

func (ur *InMemoryUserRepositoryImpl) Create(user models.User) error {
	ur.users = append(ur.users, user)
	return nil
}

func (ur *InMemoryUserRepositoryImpl) FindByUUID(uuid string) (models.User, bool) {
	for _, user := range ur.users {
		if user.UUID == uuid {
			return user, true
		}
	}

	return models.User{}, false
}

func (ur *InMemoryUserRepositoryImpl) Update(uuid string, user models.User) error {
	for index, u := range ur.users {
		if u.UUID == uuid {
			ur.users[index] = user
			return nil
		}
	}
	return fmt.Errorf("user not found")
}

func (ur *InMemoryUserRepositoryImpl) Delete(uuid string) error {
	for i, user := range ur.users {
		if user.UUID == uuid {
			ur.users = slices.Delete(ur.users, i, i+1)
			return nil
		}
	}
	return fmt.Errorf("user not found")
}

func (ur *InMemoryUserRepositoryImpl) FindByEmail(email string) (models.User, bool) {
	for _, user := range ur.users {
		if user.Email == email {
			return user, true
		}
	}
	return models.User{}, false
}
