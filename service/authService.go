package service

import (
	"core/dbhandler"
	"core/domain"
	"core/dto"
	"core/errs"
	"core/utils"
)

//	Login gets request DTO of type NewLoginRequestand check availability of fields. Sends email and password into dbhandler. Returns UserResponse dto struct and an error.
func Login(req dto.NewLoginRequest) (*dto.NewLoginResponse, *errs.AppError) {
	if req.Email == "" || req.Password == "" {
		return dto.ToULoginDto(nil, errs.Required, false), errs.NewNotFoundError("enter all credentials")
	}
	user, err := dbhandler.Login(req.Email, req.Password)

	if err != nil {
		if err.Message == errs.Wrong_em {
			return dto.ToULoginDto(nil, errs.Wrong_em, false), err
		} else if err.Message == errs.Wrong_pw {
			return dto.ToULoginDto(nil, errs.Wrong_pw, false), err
		} else if err.Message == errs.Unverified_em {
			return dto.ToULoginDto(nil, errs.Unverified_em, false), err
		} else {
			return dto.ToULoginDto(nil, errs.Deactivated_acc, false), err
		}
	}
	userResp := dto.ToULoginDto(user, "login successfull", true)
	return userResp, nil
}

/*	VerifyEmail gets email verification hash and email sent from controller. Binds them into a struct type of domain.EmailVerification. Sends it into dbhandler.
Returns DTO of type *dto.ErrResponce.
*/
func VerifyEmail(evpw, url_email string) *dto.ErrResponce {
	var response dto.ErrResponce
	verify := domain.EmailVerification{
		Email: url_email,
		Code:  evpw,
	}
	err := dbhandler.VerifyEmail(verify)
	if err != nil {
		response = dto.ErrResponce{
			Message: err.Message,
			Status:  false,
		}
	} else {
		response = dto.ErrResponce{
			Message: "Email Verified",
			Status:  true,
		}
	}
	return &response
}

//	GetActivationEmail gets DTO of dto.SendEmail and sent gets it's email field value and send it into dbhandler. Returns DTO type of *dto.ErrResponce.
func GetActivationEmail(req dto.SendEmail) *dto.ErrResponce {
	var response dto.ErrResponce
	err := dbhandler.ResendActivationEmail(req.Email)
	if err != nil {
		response = dto.ErrResponce{
			Message: err.Message,
			Status:  false,
		}
	} else {
		response = dto.ErrResponce{
			Message: "Success",
			Status:  true,
		}
	}
	return &response
}

/*	RecoverAccount gets DTO of type dto.SendEmail from controller. Check if email is in DTO request.
	Sends email in request DTO into dbhandler. Returns response DTO of type *dto.ErrResponce.
*/
func RecoverAccount(req dto.SendEmail) *dto.ErrResponce {
	var response dto.ErrResponce
	if req.Email == "" {
		return &dto.ErrResponce{
			Message: "Enter all required fields",
			Status:  false,
		}
	}
	err := dbhandler.RecoverEmail(req.Email)
	if err != nil {
		response = dto.ErrResponce{
			Message: err.Message,
			Status:  false,
		}
	} else {
		response = dto.ErrResponce{
			Message: "Please check your inbox",
			Status:  true,
		}
	}
	return &response
}

/*	ResetPassword gets DTO of type dto.PasswordReset, email verification hash and url email from controller. Check if all required fields are
	in DTO. sends values of email verification hash, url email, email in DTO, new password in DTO, confirmed password in DTO into service.
	Returns response DTO of type dto.ErrResponce.
*/
func ResetPassword(req dto.PasswordReset, evpw, url_email string) *dto.ErrResponce {
	var response dto.ErrResponce
	if req.Email == "" || req.NewPassword == "" || req.ConfirmPassword == "" {
		return &dto.ErrResponce{
			Message: "Enter all required fields",
			Status:  false,
		}
	}
	if utils.IsValidPassword(req.NewPassword) == false {
		return &dto.ErrResponce{
			Message: errs.Invalid_pw,
			Status:  false,
		}
	}
	err := dbhandler.ResetPassword(evpw, url_email, req.Email, req.NewPassword, req.ConfirmPassword)
	if err != nil {
		response = dto.ErrResponce{
			Message: err.Message,
			Status:  false,
		}
	} else {
		response = dto.ErrResponce{
			Message: "Password reset successfull",
			Status:  true,
		}
	}
	return &response
}

//	FindIfEmailExists gets email from controller and send it into dbhandler. Returns DTO type *dto.ErrResponce as response.
func FindIfEmailExists(email string) *dto.ErrResponce {
	if dbhandler.FindIfEmailExists(email) == true {
		return &dto.ErrResponce{
			Message: "A user has already exists by this email",
			Status:  false,
		}
	}
	return &dto.ErrResponce{
		Message: "Success!",
		Status:  true,
	}
}

/*	VerifyMobile gets DTO of type dto.MobileVerifyCode and user id from controller. Check if all required fields are in DTO
sends mobile number, code values in DTO and user id into dbhandler. Returns *dto.ErrResponce
*/
func VerifyMobile(req dto.MobileVerifyCode, id string) *dto.ErrResponce {
	var response dto.ErrResponce
	if req.Code == "" || req.Mobile == "" {
		return &dto.ErrResponce{
			Message: "Enter all required fields",
			Status:  false,
		}
	}
	verify := domain.MobileVerification{
		MobileNumber: req.Mobile,
		Code:         req.Code,
		UserId:       id,
	}
	err := dbhandler.VerifyMobileNo(verify)
	if err != nil {
		response = dto.ErrResponce{
			Message: err.Message,
			Status:  false,
		}
	} else {
		response = dto.ErrResponce{
			Message: "Mobile Verification Successfull",
			Status:  true,
		}
	}
	return &response
}

/*	GetMobileVerificationCode gets DTO of dto.NewMobileVerificationRequest and user id.
check if all required fields are in the request. sends mobile number field's value in DTO and user id into dbhandler.
Returns DTO type of *dto.ErrResponce.
*/
func GetMobileVerificationCode(req dto.NewMobileVerificationRequest, id string) *dto.ErrResponce {
	var response dto.ErrResponce
	if req.Mobile == "" {
		return &dto.ErrResponce{
			Message: "Enter all required fields",
			Status:  false,
		}
	}
	if utils.MobileNumberValidation(req.Mobile) == false {
		return &dto.ErrResponce{
			Message: "Your mobile number is invalid",
			Status:  false,
		}
	}
	mbReq := dto.NewMobileVerificationRequest{
		Mobile: req.Mobile,
	}
	err := dbhandler.SendMobileVerificationCode(mbReq.Mobile, id)
	if err != nil {
		response = dto.ErrResponce{
			Message: err.Message,
			Status:  false,
		}
	} else {
		response = dto.ErrResponce{
			Message: "Please enter OTP sent to your mobile number to verify",
			Status:  true,
		}
	}
	return &response
}
