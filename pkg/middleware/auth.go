package middleware

import (
	"github.com/gosquash/api/pkg/structs"

	"github.com/labstack/echo/v4"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		req := c.Request()

		// Check if the user is authenticated
		// If not, return an error
		token := req.Header.Get("Authorization")
		if len(token) == 0 {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		user := structs.GetUserByToken(token)
		if user == nil {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		c.Set("user", user)

		return next(c)
	}
}
