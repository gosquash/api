package errors

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	var message string

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	} else {
		message = err.Error()
	}

	c.JSON(code, echo.Map{
		"error": HTTPError{
			message,
			code,
		},
	})
}
