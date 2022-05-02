package main

import (
	"os"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type application struct {
	logger *log.Logger
	server *echo.Echo
}

func (app *application) setup() error {
	app.setupLog()
	app.setupServer()

	return nil
}

func (app *application) setupLog() {
	app.logger = log.New()
	app.logger.SetFormatter(&log.JSONFormatter{})
	app.logger.SetOutput(os.Stdout)
	app.logger.SetLevel(log.DebugLevel)
}

func (app *application) setupServer() {
	app.server = echo.New()
	app.routes()
}

func (app *application) startServer() error {
	err := app.server.Start(":9100")
	if err != nil {
		app.logger.Fatal(err)
		return err
	}
	return nil
}
