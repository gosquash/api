package routes

import (
	"gosquash/api/internal/routes/auth"
	"gosquash/api/internal/routes/games"
	"gosquash/api/internal/routes/groups"
	"gosquash/api/internal/routes/me"
	"gosquash/api/internal/routes/users"

	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {

	// Auth
	auth.Init(e.Group("/auth"))

	// Games
	games.Init(e.Group("/games"))

	// Groups
	groups.Init(e.Group("/groups"))

	// Me
	me.Init(e.Group("/me"))

	// Users
	users.Init(e.Group("/users"))
}
