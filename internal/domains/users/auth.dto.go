package users

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/DarrelA/starter-go-postgresql/internal/utils/err_rest"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("passwd", validatePassword)
}

type User struct {
	ID        int64      `json:"ID"`
	UUID      *uuid.UUID `json:"uuid"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Password  string     `json:"password"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type RegisterInput struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50,alpha"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50,alpha"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,passwd"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	UUID      *uuid.UUID `json:"uuid"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
}

type UserRecord struct {
	UUID      *uuid.UUID `json:"uuid"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

var validationMessages = map[string]string{
	"required": "not be empty",
	"min":      "be at least %s characters long",
	"max":      "be at most %s characters long",
	"alpha":    "contain only alphabetic characters",
	"email":    "be a valid email address",
	"passwd":   "contain at least one number, one uppercase letter, one lowercase letter, and one special character",
}

/*
For each validation error, the corresponding user-friendly message is fetched from the validationMessages map.
If a tag does not have a predefined message, a default message "be valid" is used.
The message is constructed using the field name and the user-friendly message.
*/
func ValidateStruct[T any](payload T) *err_rest.RestErr {
	var validationErrors []string

	// Pre-process the input if it is of type *RegisterInput
	if input, ok := any(&payload).(*RegisterInput); ok {
		preProcessInput(input)
	}

	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			tag := err.Tag()
			messageTemplate, exists := validationMessages[tag]

			if !exists {
				messageTemplate = "be valid"
			}

			fieldName := camelToSnakeCase(err.Field())
			message := ""

			if strings.Contains(messageTemplate, "%s") {
				message = fmt.Sprintf("the field [%s] should "+messageTemplate, fieldName, err.Param())
			} else {
				message = fmt.Sprintf("the field [%s] should %s", fieldName, messageTemplate)
			}

			validationErrors = append(validationErrors, message)
		}

		fullMessage := "validation error: " + strings.Join(validationErrors, "\n")
		return err_rest.NewBadRequestError(fullMessage)
	}

	return nil
}

func preProcessInput(payload *RegisterInput) {
	payload.FirstName = strings.TrimSpace(strings.ToLower(payload.FirstName))
	payload.LastName = strings.TrimSpace(strings.ToLower(payload.LastName))
	payload.Email = strings.TrimSpace(strings.ToLower(payload.Email))
	payload.Password = strings.TrimSpace(payload.Password)
}

func camelToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var hasNumber, hasUpper, hasLower, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasNumber && hasUpper && hasLower && hasSpecial
}
