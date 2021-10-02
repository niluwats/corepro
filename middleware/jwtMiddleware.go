package middleware

import (
	"log"
	"core/utils"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func JWTProtected() func(*fiber.Ctx) error {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config : ", err)
	}
	con := jwtware.Config{
		SigningKey: []byte(config.JwtSecretKey),
		ContextKey: "jwt",
	}
	return jwtware.New(con)
}
