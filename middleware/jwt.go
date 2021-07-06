package middleware

import (
    "github.com/gofiber/fiber/v2"
    "server/security"
    "net/http"
    "server/util"

    jwtware "github.com/gofiber/jwt/v2"
)

func AuthRequired(ctx *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:    security.JwtSecretKey,
		SigningMethod: security.JwtSigningMethod,
		TokenLookup:   "header:Authorization",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.
				Status(http.StatusUnauthorized).
				JSON(util.NewJError(err))
		},
	})(ctx)
}