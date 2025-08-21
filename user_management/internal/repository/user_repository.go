package repository

import (
	"log"

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

func (ur *InMemoryUserRepositoryImpl) FindAll() {
	log.Println("GetAllUser into InMemoryUserRepository")
}

func (ur *InMemoryUserRepositoryImpl) Create(user models.User) error {
	ur.users = append(ur.users, user)
	return nil
}

func (ur *InMemoryUserRepositoryImpl) FindByUUID() {

}

func (ur *InMemoryUserRepositoryImpl) Update() {

}

func (ur *InMemoryUserRepositoryImpl) Delete() {

}

func (ur *InMemoryUserRepositoryImpl) FindByEmail(email string) (models.User, bool) {
	for _, user := range ur.users {
		if user.Email == email {
			return user, true
		}
	}
	return models.User{}, false
}
