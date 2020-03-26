package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Root returns Status OK
func Root(c echo.Context) error {
	return c.String(http.StatusOK, "Running echo API v1")
}
