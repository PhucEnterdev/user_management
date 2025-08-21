package validation

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidation() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validator engine")
	}
	RegisterCustomValidation(v)
	return nil
}

/** function này custom lỗi trả về cho endpoint */
func HandlerValidationError(err error) gin.H {
	if validationError, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, e := range validationError {

			log.Printf("%+v", e.Error())
			log.Printf("%+v", e.Field())
			log.Printf("%+v", e.Tag())

			switch e.Tag() {
			case "gt":
				errors[e.Field()] = fmt.Sprintf("%s phải lớn hơn %s", e.Field(), e.Param())
			case "gte":
				errors[e.Field()] = fmt.Sprintf("%s phải lớn hơn hoặc bằng %s ", e.Field(), e.Param())
			case "lte":
				errors[e.Field()] = fmt.Sprintf("%s phải bé hơn hoặc bằng %s", e.Field(), e.Param())
			case "lt":
				errors[e.Field()] = fmt.Sprintf("%s phải bé hơn %s", e.Field(), e.Param())
			case "uuid":
				errors[e.Field()] = fmt.Sprintf("%s không hợp lệ", e.Field())
			case "slug":
				errors[e.Field()] = fmt.Sprintf("%s chỉ được chứa chữ thường, dấu gạch nối hoặc dấu chấm", e.Field())
			case "min":
				errors[e.Field()] = fmt.Sprintf("%s phải lớn hơn %s ký tự", e.Field(), e.Param())
			case "max":
				errors[e.Field()] = fmt.Sprintf("%s phải ít hơn %s ký tự", e.Field(), e.Param())
			case "min_int":
				errors[e.Field()] = fmt.Sprintf("%s phải có giá trị lớn hơn %s", e.Field(), e.Param())
			case "max_int":
				errors[e.Field()] = fmt.Sprintf("%s phải có gí trị nhỏ hơn %s", e.Field(), e.Param())
			case "oneof":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ", ")
				errors[e.Field()] = fmt.Sprintf("%s phải là 1 trong các giá trị %s", e.Field(), allowedValues)
			case "required":
				errors[e.Field()] = fmt.Sprintf("%s là bắt buộc", e.Field())
			case "search":
				errors[e.Field()] = fmt.Sprintf("%s là chữ cái, số và khoảng trắng", e.Field())
			case "email":
				errors[e.Field()] = fmt.Sprintf("%s không hợp lệ", e.Field())
			case "datetime":
				errors[e.Field()] = fmt.Sprintf("%s phải theo định dạng dd-mm-yyyy", e.Field())
			case "email_advanced":
				errors[e.Field()] = fmt.Sprintf("%s này nằm trong danh sách bị cấm", e.Field())
			case "password_strong":
				errors[e.Field()] = fmt.Sprintf("%s này chưa đủ mạnh (từ 8 ký tự trở lên, bao gồm chữ hoa, chữ thường, số và ký tự đặc biệt)", e.Field())
			case "file_extension":
				allowedValues := strings.Join(strings.Split(e.Param(), " "), ", ")
				errors[e.Field()] = fmt.Sprintf("%s chỉ cho phép những đuôi file %s", e.Field(), allowedValues)
			}
		}
		return gin.H{
			"error": errors,
		}
	}
	return gin.H{
		"error":  "Yêu cầu không hợp lệ ",
		"detail": err.Error(),
	}
}
