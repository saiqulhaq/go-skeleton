package middleware

import (
	"github.com/gofiber/fiber/v2"
	apperr "github.com/saiqulhaq/blog-mysql/error"
	"github.com/saiqulhaq/blog-mysql/internal/http/auth"
)

func VerifyJWTToken(c *fiber.Ctx) error {
	if err := auth.VerifyToken(c); err != nil {
		return c.Status(apperr.ErrInvalidToken().HTTPCode).JSON(apperr.ErrInvalidToken())
	}

	return c.Next()
}
