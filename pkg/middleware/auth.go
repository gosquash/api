package middleware

import (
	"log"
	"strings"

	"github.com/gosquash/api/pkg/structs"

	"github.com/labstack/echo/v4"
)

// Authentication middleware
func Authentication(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		req := c.Request()

		// Check if Authorization header is present
		token := req.Header.Get("Authorization")
		if token == "" {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		// Check if Authorization header is in the correct format
		parts := strings.SplitN(strings.TrimSpace(token), " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("Invalid Authorization Token format: %v", parts)
			return echo.NewHTTPError(401, "Unauthorized")
		}

		// Get user by token
		userToken := parts[1]
		user := structs.GetUserByToken(userToken)
		if user == nil {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		c.Set("user", user)

		return next(c)
	}
}
