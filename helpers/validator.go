package helpers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func TranslateErrorMessage(err error) map[string][]string {
	errorsMap := make(map[string][]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())

			switch fieldError.Tag() {
			case "required":
				errorsMap[field] = append(errorsMap[field], fmt.Sprintf("%s is required", field))
			case "email":
				errorsMap[field] = append(errorsMap[field], "Invalid email format")
			case "unique":
				errorsMap[field] = append(errorsMap[field], fmt.Sprintf("%s already exists", field))
			case "min":
				errorsMap[field] = append(errorsMap[field], fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param()))
			case "max":
				errorsMap[field] = append(errorsMap[field], fmt.Sprintf("%s must be at most %s characters", field, fieldError.Param()))
			case "numeric":
				errorsMap[field] = append(errorsMap[field], fmt.Sprintf("%s must be a number", field))
			default:
				errorsMap[field] = append(errorsMap[field], "Invalid value")
			}
		}
	}

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			errorsMap["email"] = append(errorsMap["email"], "Email already exists")
		} else if err == gorm.ErrRecordNotFound {
			errorsMap["error"] = append(errorsMap["error"], "Record not found")
		}
	}

	return errorsMap
}

func IsDuplicateEntryError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}
