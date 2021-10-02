package middleware

//	IsAuthenticated parses cookie and check if calling route is authenticated for the user or not. Returns an error.
// func IsAuthenticated(c *fiber.Ctx) error {
// 	cookie := c.Cookies("jwt")
// 	if _, err := utils.ParseJwt(cookie); err != nil {
// 		c.Status(fiber.StatusUnauthorized)
// 		return c.JSON(fiber.Map{
// 			"message": "missing token",
// 		})
// 	}
// 	return c.Next()
// }

//	GetIdFromCookie returns user Id from cookie as a string
// func GetIdFromCookie(c *fiber.Ctx) string {
// 	cookie := c.Cookies("jwt")
// 	Id, _ := utils.ParseJwt(cookie)
// 	return Id
// }
