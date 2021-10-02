package routes

import (
	"core/controller"
	"core/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUp(app *fiber.App) {
	//app.Get("/", controller.RedirectTo)
	app.Post("/v1/auth/users", controller.NewUser)
	app.Post("/v1/auth/login", controller.Login)
	app.Post("/v1/auth/verifyemail/:email/:evpw", controller.VerifyEmail)

	app.Post("/v1/auth/sendactivationemail", controller.SendActivationEmail)
	app.Post("/v1/auth/recoveraccount", controller.RecoverAccount)
	app.Put("/v1/auth/resetpassword/:email/:evpw", controller.ResetPassword)
	app.Get("/v1/auth/email", controller.FindEmail)

	app.Post("/v1/auth/users/sendsms/:id", middleware.JWTProtected(), controller.GetMbVerificationCode)
	app.Post("/v1/auth/users/verifymobile/:id", middleware.JWTProtected(), controller.VerifyMobile)

}
