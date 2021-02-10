package routers

import (
	"log"
	"main/core"
	"main/core/business"
	"main/core/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserRouter struct {
	Name string
	g    *echo.Group
}

func (r *UserRouter) Connect(s *core.Server) {
	r.g = s.Echo.Group(r.Name)

	user := business.UserBusiness{
		DB: s.DB,
	}

	customAuth, err := business.NewAuthBusiness(s.DB, s.Config.AuthSecret)
	if err != nil {
		log.Fatal("Failed to authorizing")
	}
	r.g.GET("/", func(c echo.Context) error {
		user := c.Get("user")
		return c.JSON(http.StatusOK, user)
	}, s.AuthWiddlewareJWT.Auth)

	r.g.POST("/", func(c echo.Context) (err error) {
		profile := new(models.UserProfile)

		if err = c.Bind(profile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(profile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		err = user.Create(*profile)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, profile)
	})

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

	// [PUT] Input: Body (UserProfileUpdated)

	// [DELETE] Input user's id

}
