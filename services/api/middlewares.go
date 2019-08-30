package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/crazyfacka/yanbapp-api/domain/users"
	"github.com/crazyfacka/yanbapp-api/tools"
)

var sessionException = []string{
	"/session",
}

// CheckSession validates if the user has a valid session
func CheckSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		for _, item := range sessionException {
			if c.Request().RequestURI == item {
				return next(c)
			}
		}

		cookie, err := c.Cookie(users.SessionCookieName)
		if err != nil {
			return c.String(http.StatusUnauthorized, tools.JSONMessage("valid cookie not found"))
		}

		u := &users.User{
			SessionKey: cookie.Value,
		}

		if err = cacheRepository.LoadSession(u); err != nil {
			return c.String(http.StatusUnauthorized, tools.JSONMessage("valid cookie not found"))
		}

		c.Set("user", u)

		return next(c)
	}
}
