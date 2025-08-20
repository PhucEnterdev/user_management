package handler

import (
	"net/http"

	"enterdev.com.vn/user_management/internal/models"
	"enterdev.com.vn/user_management/internal/services"
	"enterdev.com.vn/user_management/internal/validation"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {

}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, validation.HandlerValidationError(err))
	}
	uh.service.CreateUser(user)
}

func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {

}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {

}

func (uh *UserHandler) DeleteUser(ctx *gin.Context) {

}
