package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/krish-srivastava-2305/internals/services"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.JSON(401, echo.Map{"error": "Unauthorized"})
		}

		token := cookie.Value
		jwtToken, err := services.ValidateToken(token)

		if err != nil {
			return c.JSON(401, echo.Map{"error": "Unauthorized"})
		}

		claims := jwtToken.Claims.(jwt.MapClaims)
		email := claims["email"].(string)

		c.Set("email", email)

		return next(c)
	}
}
