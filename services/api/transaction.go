package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/crazyfacka/yanbapp-api/domain/transactions"
	"github.com/crazyfacka/yanbapp-api/domain/users"
	"github.com/crazyfacka/yanbapp-api/tools"
)

func getTransactions(c echo.Context) error {
	var ts []transactions.Transaction

	u := c.Get("user").(*users.User)
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.String(http.StatusBadRequest, tools.JSONMessage(err.Error()))
	}

	if ts, err = dbRepository.GetTransactions(u, uint(id)); err != nil {
		return c.String(http.StatusInternalServerError, tools.JSONMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, ts)
}

func saveTransaction(c echo.Context) error {
	var (
		insertID int64
		repeatID int64
		err      error
	)

	t := new(transactions.Transaction)

	u := c.Get("user").(*users.User)

	if err = c.Bind(t); err != nil {
		return c.String(http.StatusBadRequest, tools.JSONMessage(err.Error()))
	}

	if insertID, repeatID, err = dbRepository.SaveTransaction(u, t); err != nil {
		return c.String(http.StatusBadRequest, tools.JSONMessage(err.Error()))
	}

	rsp := map[string]interface{}{
		"id":        insertID,
		"repeat_id": repeatID,
	}

	return c.String(http.StatusOK, tools.JSONMessageMap("OK", rsp))
}
