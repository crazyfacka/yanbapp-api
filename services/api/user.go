package api

import (
	"net/http"

	"github.com/crazyfacka/yanbapp-api/tools"

	"github.com/labstack/echo"
	"github.com/crazyfacka/yanbapp-api/domain/users"
)

func createSession(c echo.Context) error {
	u := new(users.User)

	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, tools.JSONMessage(err.Error()))
	}

	if err := dbRepository.ValidateUser(u); err != nil {
		return c.String(http.StatusBadRequest, tools.JSONMessage(err.Error()))
	}

	if err := cacheRepository.CreateSession(u); err != nil {
		return c.String(http.StatusBadRequest, tools.JSONMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, u)
}
