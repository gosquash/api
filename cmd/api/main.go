package main

import (
	"os"

	"github.com/gosquash/api/internal/db"
	"github.com/gosquash/api/internal/routes"
	"github.com/gosquash/api/internal/validator"
	"github.com/gosquash/api/pkg/errors"
	"github.com/gosquash/api/pkg/structs"

	"github.com/labstack/echo/v4"
)

func main() {

	// Initialize database
	db.Init()

	// Sync database
	db.DB.AutoMigrate(
		&structs.User{},
		&structs.Game{},
		&structs.Group{},
		&structs.GroupMember{},
		&structs.Player{},
	)

	e := echo.New()

	// Validator
	e.Validator = validator.NewValidator()

	// Custom error handler.
	e.HTTPErrorHandler = errors.ErrorHandler

	// Define routes
	routes.Init(e)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
