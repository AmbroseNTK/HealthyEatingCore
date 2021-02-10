package middlewares

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
)

type AuthMiddleware struct {
	signer *jwt.HMAC
}

func NewAuthMiddleware(key string) *AuthMiddleware {
	signer, err := jwt.NewHMAC(jwt.HS512, []byte(key))
	if err != nil {
		log.Fatalln("Cannot create new JWT Signer")
	}
	return &AuthMiddleware{
		signer: signer,
	}
}

func (m *AuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		idToken := c.Request().Header.Get("Authorization")
		if idToken == "" {
			idToken = c.QueryParam("token")
		}
		if idToken == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing ID Token")
		}
		token, tokenError := m.signer.Check([]byte(idToken))
		if tokenError != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, tokenError.Error())
		}
		c.Set("user", token.Set)
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}
