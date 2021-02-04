package routers

import (
	"main/core"
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
	r.g.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello, world",
		})
	})

	r.g.POST("/", func(c echo.Context) (err error) {
		profile := new(models.UserProfile)

		if err = c.Bind(profile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(profile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, profile)
	})

}
