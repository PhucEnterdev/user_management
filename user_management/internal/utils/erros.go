package utils

type AppErrror struct {
	Message string
	Code    string
	Error   error
}

func NewError(message string, code string) *AppErrror {
	return &AppErrror{
		Message: message,
		Code:    code,
		Error:   ,
	}
}
