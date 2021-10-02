package service

import (
	"core/dbhandler"
	"core/domain"
	"core/dto"
	"core/errs"
	"core/utils"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser gets CreateUser dto request struct and check availabiblity of required fields, validity of entered fields, encrypt password. Binds request data into a domain.User struct. Sends user struct to CreateUser in dbhandler. Returns *NewUserResponse dto as a struct and an error of type *errs.AppError.
func CreateUser(req dto.NewUserRequest) (*dto.NewUserResponse, *errs.AppError) {
	if req.Name == "" || req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		return dto.ToUserDto(nil, errs.Required, false), errs.NewUnexpectedError("enter all required fields")
	}
	if req.ConfirmPassword != req.Password {
		return dto.ToUserDto(nil, errs.Wrong_pw, false), errs.NewUnexpectedError("miss matched passwords")
	}
	if utils.IsValidEmail(req.Email) == false {
		return dto.ToUserDto(nil, errs.Wrong_pw, false), errs.NewUnexpectedError("invalid email address")
	}
	if utils.IsValidPassword(req.Password) == false {
		return dto.ToUserDto(nil, errs.Invalid_pw, false), errs.NewUnexpectedError("Your password should be more than 6 characters long with including atleast one special character,capital letter and a number !")
	}

	hashedPw, err := bcrypt.GenerateFromPassword([]byte(req.ConfirmPassword), bcrypt.DefaultCost)
	if err != nil {
		return dto.ToUserDto(nil, err.Error(), false), errs.NewUnexpectedError(err.Error())
	}

	user := domain.User{
		Name:           req.Name,
		Email:          req.Email,
		Password:       string(hashedPw),
		EmailVerified:  false,
		MobileVerified: false,
		Activated:      false,
	}

	resp, err1 := dbhandler.CreateUser(user)
	if err1 != nil {
		return dto.ToUserDto(nil, err1.Message, false), err1
	}
	return dto.ToUserDto(resp, "User has been created, Please check your inbox to verify your email", true), nil
}
