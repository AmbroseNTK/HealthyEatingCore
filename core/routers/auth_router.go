package routers

import (
	"log"
	"main/core"
	"main/core/business"
	"main/core/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthRouter struct {
	Name string
	g    *echo.Group
}

func (r *AuthRouter) Connect(s *core.Server) {

	r.g = s.Echo.Group(r.Name)

	customAuth, err := business.NewAuthBusiness(s.DB, s.Config.AuthSecret)
	if err != nil {
		log.Fatal("Failed to authorizing")
	}
	r.g.POST("/register", func(c echo.Context) (err error) {
		userAuth := new(models.UserAuth)
		if err = c.Bind(userAuth); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(userAuth); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = customAuth.Register(userAuth)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return c.String(http.StatusOK, "Registrated")
	})

	r.g.POST("/login", func(c echo.Context) (err error) {
		userAuth := new(models.UserAuth)
		if err = c.Bind(userAuth); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(userAuth); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		token, loginErr := customAuth.Login(userAuth)
		if loginErr != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, loginErr.Error())
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"token": token,
		})
	})
}
