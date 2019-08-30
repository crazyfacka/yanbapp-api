package api

import (
	"net/http"
	"strconv"

	"github.com/crazyfacka/yanbapp-api/domain/accounts"

	"github.com/crazyfacka/yanbapp-api/tools"

	"github.com/labstack/echo"
	"github.com/crazyfacka/yanbapp-api/domain/users"
)

func getAccounts(c echo.Context) error {
	u := c.Get("user").(*users.User)

	if err := dbRepository.GetAccounts(u); err != nil {
		return c.String(http.StatusInternalServerError, tools.JSONMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, u.Accounts)
}

func saveAccount(c echo.Context) error {
	var (
		insertID int64
		err      error
	)

	a := new(accounts.Account)

	u := c.Get("user").(*users.User)

	if err = c.Bind(a); err != nil {
		return c.String(http.StatusBadRequest, tools.JSONMessage(err.Error()))
	}

	if insertID, err = dbRepository.SaveAccount(u, a); err != nil {
		return c.String(http.StatusBadRequest, tools.JSONMessage(err.Error()))
	}

	return c.String(http.StatusOK, tools.JSONMessage(strconv.Itoa(int(insertID))))
}
