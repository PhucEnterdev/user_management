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

type GetUserParams struct {
	Search string `form:"search" binding:"omitempty,min=3,max=50,search"`
	Page   int    `form:"page" binding:"omitempty,gte=1,lte=100"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {
	var params GetUserParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	users, err := uh.service.GetAllUser(params.Search, params.Page, params.Limit)
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
	var params GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.ResponseValidator(ctx, err)
		return
	}
	updatedUser, err := uh.service.UpdateUser(params.UUID, user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapUserToDTO(updatedUser)
	utils.ResponseSuccess(ctx, http.StatusOK, &userDTO)
}

func (uh *UserHandler) DeleteUser(ctx *gin.Context) {

}
