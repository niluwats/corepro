package controller

import (
	"core/dto"
	"core/service"

	"github.com/gofiber/fiber/v2"
)

//	VerifyEmail gets url params email verification hash and email. Sends them into service. Returns http status, status, message for the response.
func VerifyEmail(c *fiber.Ctx) error {
	evpw := c.Params("evpw")
	url_email := c.Params("email")
	appErr := service.VerifyEmail(evpw, url_email)
	if appErr.Status != true {
		getStatus(c, 500, appErr)
	} else {
		getStatus(c, 200, appErr)
	}
	return nil
}

//	SendActivationEmail gets request and parses request body into a struct. Sends it into service. Returns http status and message as response.
func SendActivationEmail(c *fiber.Ctx) error {
	var req dto.SendEmail
	if err := c.BodyParser(&req); err != nil {
		getStatus(c, 422, err)
	} else {
		appErr := service.GetActivationEmail(req)
		if appErr.Status != true {
			getStatus(c, 500, appErr)
		} else {
			getStatus(c, 200, appErr)
		}
	}
	return nil
}

//	RecoverAccount gets request and parse request body into a struct to get user entered email. Send it to service. Returns http status, status and a message.
func RecoverAccount(c *fiber.Ctx) error {
	// if middleware.GetIdFromCookie(c) != "" {
	// 	getStatus(c, 401, "already have logged in !")
	// }
	var request dto.SendEmail
	if err := c.BodyParser(&request); err != nil {
		getStatus(c, 422, err)
	} else {
		appErr := service.RecoverAccount(request)
		if appErr.Status != true {
			getStatus(c, 500, appErr)
		} else {
			getStatus(c, 200, appErr)
		}
	}
	return nil
}

//	ResetPassword gets request and parse request body into a struct to get data in request. Gets url params from url. Pass them into service.Returns http status, status and a message.
func ResetPassword(c *fiber.Ctx) error {
	var req dto.PasswordReset
	evpw := c.Params("evpw")
	url_email := c.Params("email")

	if err := c.BodyParser(&req); err != nil {
		getStatus(c, 422, err)
	} else {
		appErr := service.ResetPassword(req, evpw, url_email)
		if appErr.Status != true {
			getStatus(c, 500, appErr)
		} else {
			getStatus(c, 200, appErr)
		}
	}
	return nil
}

//	FindEmail gets request and parse data in request body and get user entered email. Send it to service. Return http status and response.
func FindEmail(c *fiber.Ctx) error {
	var req dto.SendEmail
	if err := c.BodyParser(&req); err != nil {
		getStatus(c, 422, err)
	} else {
		resp := service.FindIfEmailExists(req.Email)
		if resp.Status != true {
			getStatus(c, 500, resp)
		} else {
			getStatus(c, 200, resp)
		}
	}
	return nil
}

//	GetMbVerificationCode gets url param (userid) and gets request and parse request body into a struct. Pass them into service. Returns http status, status, message for the response.
func GetMbVerificationCode(c *fiber.Ctx) error {
	// id := middleware.GetIdFromCookie(c)
	id := c.Params("id")
	var request dto.NewMobileVerificationRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	} else {
		appErr := service.GetMobileVerificationCode(request, id)
		if appErr.Status != true {
			getStatus(c, 500, appErr)
		} else {
			getStatus(c, 200, appErr)
		}
	}
	return nil
}

//	VerifyMobile gets url param (userid) and gets request and parse request body into a struct. Pass them to service. Returns http status, status, message for the response.
func VerifyMobile(c *fiber.Ctx) error {
	// id := middleware.GetIdFromCookie(c)
	id := c.Params("id")
	var request dto.MobileVerifyCode
	if err := c.BodyParser(&request); err != nil {
		getStatus(c, 422, err)
	} else {
		appErr := service.VerifyMobile(request, id)
		if appErr.Status != true {
			getStatus(c, 500, appErr)
		} else {
			getStatus(c, 200, appErr)
		}
	}
	return nil
}
