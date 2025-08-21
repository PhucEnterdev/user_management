package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	ErrCodeBadRequest ErrorCode = "BAD_REQUEST"
	ErrCodeNotFound   ErrorCode = "NOT_FOUND"
	ErrCodeConflict   ErrorCode = "CONFLICT"
	ErrCodeInternal   ErrorCode = "INTERNAL_SERVER_ERROR"
)

type AppErrror struct {
	Message string
	Code    ErrorCode
	Err     error
}

func (ae *AppErrror) Error() string {
	return ""
}

func NewError(message string, code ErrorCode) error {
	return &AppErrror{
		Message: message,
		Code:    code,
	}
}

func WrapError(err error, msg string, code ErrorCode) error {
	return &AppErrror{
		Err:     err,
		Message: msg,
		Code:    code,
	}
}

func ResponseError(ctx *gin.Context, err error) {
	if appError, ok := err.(*AppErrror); ok {
		statusCode := httpStatusFormCode(appError.Code)
		response := gin.H{
			"error": appError.Message,
			"code":  appError.Code,
		}
		if appError.Err != nil {
			response["detail"] = appError.Err.Error()
		}
		ctx.JSON(statusCode, response)
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code":  ErrCodeInternal,
	})
}

func ResponseSuccess(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, gin.H{
		"status": "success",
		"data":   data,
	})
}

func httpStatusFormCode(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
