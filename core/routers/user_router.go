package routers

import (
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

	user.CreateIndexes()

	r.g.GET("/", func(c echo.Context) error {
		authUser := c.Get("user")
		return c.JSON(http.StatusOK, authUser)
	}, s.AuthWiddlewareJWT.Auth)

	r.g.POST("/", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		profile := new(models.UserProfile)

		if err = c.Bind(profile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(profile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if profile.Email != authUser["email"] {
			return echo.NewHTTPError(http.StatusBadRequest, "Input email and authorized email did not match")
		}

		err = user.Create(*profile)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, profile)
	}, s.AuthWiddlewareJWT.Auth)

	// [PUT] Input: Body (UserProfileUpdated)

	r.g.PUT("/", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		updatedProfile := new(models.UpdatedUserProfile)
		if err = c.Bind(updatedProfile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(updatedProfile); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = user.Update(authUser["email"].(string), *updatedProfile)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth)

	// [DELETE] Input user's id

	r.g.DELETE("/", func(c echo.Context) (err error) {
		authUser := c.Get("user").(map[string]interface{})
		err = user.Delete(authUser["email"].(string))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}, s.AuthWiddlewareJWT.Auth)

}
