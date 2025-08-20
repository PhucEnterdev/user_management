package validation

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidation(v *validator.Validate) {
	var slugRegex = regexp.MustCompile(`^[a-z0-9]+(?:[-.][a-z0-9]+)*$`)
	v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slugRegex.MatchString(fl.Field().String())
	})

	var searchRegex = regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)
	v.RegisterValidation("search", func(fl validator.FieldLevel) bool {
		return searchRegex.MatchString(fl.Field().String())
	})

	v.RegisterValidation("min_int", func(fl validator.FieldLevel) bool {
		minStr := fl.Param()
		// Base
		// 10: hệ thập phân (decimal)
		// 16: hệ thập lục phân (hex, ví dụ: "FF" => 255)
		// 2: hệ nhị phân (binary)
		// 64 ở đây là 64 bit
		minValue, err := strconv.ParseInt(minStr, 10, 64)
		if err != nil {
			return false
		}
		return fl.Field().Int() >= minValue
	})

	v.RegisterValidation("max_int", func(fl validator.FieldLevel) bool {
		maxStr := fl.Param()
		// Base
		// 10: hệ thập phân (decimal)
		// 16: hệ thập lục phân (hex, ví dụ: "FF" => 255)
		// 2: hệ nhị phân (binary)
		// 64 ở đây là 64 bit
		maxValue, err := strconv.ParseInt(maxStr, 10, 64)
		if err != nil {
			return false
		}
		return fl.Field().Int() >= maxValue
	})

	/** Kiểm tra các đuôi file nào được phép gửi lên server*/
	v.RegisterValidation("file_extension", func(fl validator.FieldLevel) bool {
		fileName := fl.Field().String()
		allowedStr := fl.Param()
		if allowedStr == "" {
			return false
		}
		/** Cắt chuỗi thành mảng các field là các đuôi file được phép gửi*/
		allowedExtension := strings.Fields(allowedStr)
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), ".")

		for _, allowed := range allowedExtension {
			if ext == strings.ToLower(allowed) {
				return true
			}
		}
		return false
	})
}
