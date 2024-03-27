package users

import (
	"gosquash/api/pkg/structs"

	"github.com/labstack/echo/v4"
)

func Init(g *echo.Group) {

	g.GET("/:id", getUserById)
}

// Create a user.
func getUserById(c echo.Context) error {

	user := structs.GetUser(c.Param("id"))

	return c.JSON(200, user)
}
