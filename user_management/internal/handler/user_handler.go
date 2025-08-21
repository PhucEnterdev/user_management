package handler

import (
	"net/http"

	"enterdev.com.vn/user_management/internal/dto"
	"enterdev.com.vn/user_management/internal/models"
	"enterdev.com.vn/user_management/internal/services"
	"enterdev.com.vn/user_management/internal/utils"
	"enterdev.com.vn/user_management/internal/validation"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

type GetUserByUUIDParam struct {
	UUID string `uri:"uuid" binding:"uuid"`
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {
	users, err := uh.service.GetAllUser()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, validation.HandlerValidationError(err))
		return
	}

	userDTO := dto.MapUsersToDTO(users)

	utils.ResponseSuccess(ctx, http.StatusOK, &userDTO)
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.ResponseValidator(ctx, err)
		return
	}
	createdUser, err := uh.service.CreateUser(user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := dto.MapUserToDTO(createdUser)

	utils.ResponseSuccess(ctx, http.StatusCreated, &userDTO)
}

func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {
	var params GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	user, err := uh.service.GetUserByUUID(params.UUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapUserToDTO(user)

	utils.ResponseSuccess(ctx, http.StatusOK, &userDTO)
}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {

}

func (uh *UserHandler) DeleteUser(ctx *gin.Context) {

}
