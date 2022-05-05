package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"pizza-waiting-time/internal/notification"
	"pizza-waiting-time/internal/organizer"
	http2 "pizza-waiting-time/internal/repository/http"
	"pizza-waiting-time/internal/repository/memory"
)

func (app application) routes() {

	repo := memory.NewPizzaTime(app.logger)

	notifications := notification.Notification{
		Logger:     app.logger,
		HTTPPoster: http2.NewClient(),
	}

	handler := organizer.NewHandler(app.logger)
	handler.PizzaTimeProvider = repo
	handler.NotificationProvider = &notifications

	group := app.server.Group("/pizza")

	group.GET("/health-check", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	})
	group.POST("/add", handler.Handle)
}
