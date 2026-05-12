package middleware

import (
	"strings"

	jwtPkg "go-trial/pkg/jwt"
	"go-trial/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtManager *jwtPkg.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Missing authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		tokenString := parts[1]

		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		if claims.Type != jwtPkg.AccessToken {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid token type")
		}

		c.Locals("claims", claims)

		return c.Next()
	}
}