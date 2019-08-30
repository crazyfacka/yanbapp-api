package api

import (
	"net/http"
	"strconv"

	"github.com/crazyfacka/yanbapp-api/domain/budget"

	"github.com/labstack/echo"
	"github.com/crazyfacka/yanbapp-api/domain/users"
	"github.com/crazyfacka/yanbapp-api/tools"
)

func getBudgetInfo(c echo.Context) error {
	u := c.Get("user").(*users.User)

	if err := dbRepository.GetBudgetInfo(u); err != nil {
		return c.String(http.StatusInternalServerError, tools.JSONMessage(err.Error()))
	}

	return c.JSON(http.StatusOK, u.BudgetGroups)
}

func saveBudgetItem(c echo.Context) error {
	var (
		insertID int64
		err      error
	)

	bi := new(budget.Item)

	u := c.Get("user").(*users.User)

	if err = c.Bind(bi); err != nil {
		return c.String(http.StatusBadRequest, tools.JSONMessage(err.Error()))
	}

	if insertID, err = dbRepository.SaveBudgetItem(u, bi); err != nil {
		return c.String(http.StatusBadRequest, tools.JSONMessage(err.Error()))
	}

	return c.String(http.StatusOK, tools.JSONMessage(strconv.Itoa(int(insertID))))
}
