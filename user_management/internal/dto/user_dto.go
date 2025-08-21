package dto

import "enterdev.com.vn/user_management/internal/models"

type UserDTO struct {
	UUID   string `json:"uuid"`
	Name   string `json:"full_name""`
	Email  string `json:"email"`
	Age    int    `json:"age"`
	Status string `json:"status"`
	Level  string `json:"level"`
}

func MapUserToDTO(user models.User) *UserDTO {
	return &UserDTO{
		UUID:   user.UUID,
		Name:   user.Name,
		Email:  user.Email,
		Age:    user.Age,
		Status: mapStatusText(user.Status),
		Level:  mapLevelText(user.Level),
	}
}

func MapUsersToDTO(users []models.User) *[]UserDTO {
	dtos := make([]UserDTO, 0, len(users))
	for _, user := range users {
		dto := UserDTO{
			UUID:   user.UUID,
			Name:   user.Name,
			Email:  user.Email,
			Age:    user.Age,
			Status: mapStatusText(user.Status),
			Level:  mapLevelText(user.Level),
		}

		dtos = append(dtos, dto)
	}
	return &dtos
}

func mapStatusText(status int) string {
	switch status {
	case 1:
		return "Show"
	case 2:
		return "Hide"
	default:
		return "None"
	}
}

func mapLevelText(level int) string {
	switch level {
	case 1:
		return "Admin"
	case 2:
		return "Member"
	default:
		return "None"
	}
}
