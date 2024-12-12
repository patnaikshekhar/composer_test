package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, struct{ Result bool }{Result: true})
}
