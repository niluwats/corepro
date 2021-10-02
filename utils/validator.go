package utils

import (
	"net/mail"
	"regexp"
	"unicode"
)

//	IsValidEmail returns if given email is a valid email or not
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

//	MobileNumberValidation returns if given mobile number is a valid mobile number or not
func MobileNumberValidation(tp string) bool {
	re := regexp.MustCompile(`^(?:7|0|(?:\+94))[0-9]{9,10}$`)
	return re.MatchString(tp)
}

//	IsValidPassword returns if given password is a valid password according to the given requirements or not
func IsValidPassword(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
