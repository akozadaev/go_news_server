package middleware

import (
	"fmt"
	jwtMiddleware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTProtected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/contrib/jwt
func JWTProtected(secretKey string) func(*fiber.Ctx) error {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)

	claims["sub"] = "5"
	claims["name"] = "akozadaev"

	token.Claims = claims
	signature := []byte(secretKey)
	fmt.Println("signature : ", signature)
	tokenString, err := token.SignedString(signature)
	fmt.Println(tokenString)
	fmt.Println(err)

	config := jwtMiddleware.Config{
		SigningKey:   jwtMiddleware.SigningKey{Key: signature},
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}
