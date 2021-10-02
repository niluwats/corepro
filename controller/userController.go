package controller

import (
	"core/dto"
	"core/service"
	"core/utils"

	"github.com/gofiber/fiber/v2"
)

//	NewUser Get user signup request and parse it.Sends parsed request as a struct to CreateUser in service. Returns http status, for the response.
func NewUser(c *fiber.Ctx) error {
	var request dto.NewUserRequest
	if err := c.BodyParser(&request); err != nil {
		getStatus(c, 422, err)
	} else {
		user, appErr := service.CreateUser(request)
		if appErr != nil {
			getStatus(c, 500, user)
		} else {
			getStatus(c, 201, &user)
		}
	}
	return nil
}

//	Login gets request and parse request body into a struct to get data in request. Pass them into service. Generates a JWT token if no error is thrown from Login service. Returns http status, status, message, user response.
func Login(c *fiber.Ctx) error {
	var request dto.NewLoginRequest
	if err := c.BodyParser(&request); err != nil {
		getStatus(c, 422, err)
	} else {
		user, appErr := service.Login(request)
		if appErr != nil {
			getStatus(c, 500, user)
		} else {
			token, err := utils.GenerateJwt(user.UserId)
			user.AccessToken = token
			if err != nil {
				getStatus(c, 401, err)
			}
			// cookie := fiber.Cookie{
			// 	Name:     "jwt",
			// 	Value:    token,
			// 	Expires:  time.Now().Add(time.Hour * 24),
			// 	HTTPOnly: true,
			// }
			// c.Cookie(&cookie)
			getStatus(c, 200, &user)
		}
	}
	return nil
}

func getStatus(c *fiber.Ctx, statusCode int, data interface{}) error {
	c.Status(statusCode)
	return c.JSON(data)
}
