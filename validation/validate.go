package validation

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("strmin", strMin)
	validate.RegisterValidation("strmax", strMax)
	validate.RegisterValidation("username", isUsername)
	validate.RegisterValidation("password", isPassword)
	validate.RegisterValidation("alphanumspecial", isAlphaNumericSpecial)
	return validate
}

func DescribeErrors(errs validator.ValidationErrors) error {
	msg := bytes.NewBufferString("")
	for _, err := range errs {
		msg.WriteString("Validation failed for '" + err.Namespace() + "': ")
		switch err.Tag() {
		case "strmin":
			msg.WriteString(" Length must be at least '" + err.Param() + "' characters.")
		case "strmax":
			msg.WriteString(" Must be maxium length of '" + err.Param() + "'.")
		case "username":
			msg.WriteString(" Must be 3 to 20 characters containing any lowercase letter, digit, underscore or hyphen.")
		case "password":
			msg.WriteString(` Must be 6 to 30 characters containing at least one digit, one uppercase letter, one lowercase letter and one special symbol (“@#$%”).`)
		case "alphanumspecial":
			msg.WriteString(` Zero or more characters containing any uppercase/lowercase letter, digit, underscore, hyphen, space or special symbol("@#$%-") only.`)
		}
		msg.WriteString("\n")
	}
	return errors.New(msg.String())
}

func strMin(fl validator.FieldLevel) bool {
	// Trim the string value first
	value := strings.TrimSpace(fl.Field().String())
	p, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}
	if len(value) < p {
		return false
	}
	return true
}

func strMax(fl validator.FieldLevel) bool {
	// If string is empty, skip the validation
	if fl.Field().String() == "" {
		return true
	}
	// Trim the string value first
	value := strings.TrimSpace(fl.Field().String())
	p, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}
	if len(value) > p {
		return false
	}
	return true
}

// 3 to 20 characters containing any lowercase letter, digit, underscore or hyphen
func isUsername(fl validator.FieldLevel) bool {
	isValid := regexp.MustCompile(`^[\da-z_-]{3,20}$`).MatchString
	return isValid(fl.Field().String())
}

// Whole combination is means, 6 to 30 characters string with at least one digit, one upper case letter, one lower case letter and one special symbol (“@#$%-”).
// This regular expression pattern is very useful to implement a strong and complex password.
func isPassword(fl validator.FieldLevel) bool {
	// isValid := regexp.MustCompile(`^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*[@#$%]).{6,20}$`).MatchString
	value := fl.Field().String()
	// Must contain at least one digit
	if isMatch := regexp.MustCompile(`[\d]+`).MatchString(value); !isMatch {
		return false
	}
	// Must contain at least one lowercase letter
	if isMatch := regexp.MustCompile(`[a-z]+`).MatchString(value); !isMatch {
		return false
	}
	// Must contain at least one uppercase letter
	if isMatch := regexp.MustCompile(`[A-Z]+`).MatchString(value); !isMatch {
		return false
	}
	// Must contain at least one special character
	if isMatch := regexp.MustCompile(`[@#$%-]+`).MatchString(value); !isMatch {
		return false
	}
	// Must be 6-30 characters
	if isMatch := regexp.MustCompile(`^[\w\s@#$%-]{6,30}$`).MatchString(value); !isMatch {
		return false
	}
	return true
}

// Zero or more characters containing any uppercase/lowercase letter, digit, underscore, hyphen, space or special symbol("@#$%").
func isAlphaNumericSpecial(fl validator.FieldLevel) bool {
	isValid := regexp.MustCompile(`^[\w @#$%-]*$`).MatchString
	return isValid(fl.Field().String())
}
