package main

import (
	"fmt"
	"log"
	"os"

	"gosquash/api/internal/db"
	"gosquash/api/internal/routes"
	"gosquash/api/internal/validator"
	"gosquash/api/pkg/errors"
	"gosquash/api/pkg/structs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {

	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db.Init()

	// Sync database
	err = db.DB.AutoMigrate(
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

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
