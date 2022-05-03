package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app application) routes() {
	group := app.server.Group("/pizza")

	group.GET("/health-check", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	})
	group.POST("/add", nil)
}
