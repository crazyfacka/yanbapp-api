package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/crazyfacka/yanbapp-api/repositories/cache"
	"github.com/crazyfacka/yanbapp-api/repositories/db"
)

var (
	e               *echo.Echo
	dbRepository    *db.DB
	cacheRepository *cache.Redis
)

// Start creates a new instance of the server
func Start(port int, dbRepositoryArg *db.DB, cacheRepositoryArg *cache.Redis) {
	dbRepository = dbRepositoryArg
	cacheRepository = cacheRepositoryArg

	e = echo.New()
	e.Use(middleware.Logger())
	e.Use(CheckSession)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "These are not the droids you're looking for...")
	})

	e.POST("/session", createSession)

	e.GET("/budget", getBudgetInfo)
	e.PUT("/budget/:id", saveBudgetItem)

	e.GET("/accounts", getAccounts)
	e.PUT("/accounts/:id", saveAccount)
	e.POST("/accounts", saveAccount)

	e.GET("/transactions/:id", getTransactions)
	e.PUT("/transactions/:id", saveTransaction)
	e.POST("/transactions", saveTransaction)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(port)))
}

// Stop gracefully terminates the API server
func Stop(ctx context.Context) {
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
