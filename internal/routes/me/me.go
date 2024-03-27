package me

import (
	"gosquash/api/internal/db"
	"gosquash/api/pkg/middleware"
	"gosquash/api/pkg/structs"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Init(g *echo.Group) {
	g.Use(middleware.Authentication)

	g.GET("", getMe)
	g.PATCH("", updateMe)
}

func getMe(c echo.Context) error {

	user := c.Get("user")

	return c.JSON(200, user)
}

func updateMe(c echo.Context) error {

	user := c.Get("user").(*structs.User)

	var body struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
	}

	if err := c.Bind(&body); err != nil {
		return err
	}

	if result := db.DB.Model(&user).Updates(&body).Find(&user); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Error updating user",
		})
	}

	return c.String(http.StatusAccepted, "")
}