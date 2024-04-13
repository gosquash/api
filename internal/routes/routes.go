package routes

import (
	"github.com/gosquash/api/internal/routes/auth"
	"github.com/gosquash/api/internal/routes/games"
	"github.com/gosquash/api/internal/routes/groups"
	"github.com/gosquash/api/internal/routes/me"
	"github.com/gosquash/api/internal/routes/users"

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
