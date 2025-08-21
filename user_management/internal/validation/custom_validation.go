package validation

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"enterdev.com.vn/user_management/internal/utils"
	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidation(v *validator.Validate) {
	var blockedDomains = map[string]bool{
		"blacklist.com": true,
		"edu.vn":        true,
		"hacker.com":    true,
	}

	// block email không được phép đăng ký
	v.RegisterValidation("email_advanced", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()

		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			return false
		}

		domain := utils.NormalizeString(parts[1])
		return !blockedDomains[domain]
	})

	// đặt password mạnh
	v.RegisterValidation("password_strong", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()

		if len(password) < 8 {
			return false
		}

		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>?/\\|]`).MatchString(password)

		return hasLower && hasUpper && hasDigit && hasSpecial
	})

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
