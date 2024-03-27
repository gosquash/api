package main

import (
	"gosquash/api/internal/db"
	"gosquash/api/internal/routes"
	"gosquash/api/pkg/structs"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError

	he, ok := err.(*echo.HTTPError)
	if ok {
		code = he.Code
	}

	c.JSON(code, echo.Map{
		"error": echo.Map{
			"message": he.Message.(string),
			"status":  code,
		},
	})
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {

	// sync database
	db.DB.AutoMigrate(
		&structs.User{},
		&structs.Game{},
		&structs.Player{},
	)

	e := echo.New()

	// Validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Custom error handler.
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Define routes
	routes.Init(e)

	e.Logger.Fatal(e.Start(":1323"))
}
