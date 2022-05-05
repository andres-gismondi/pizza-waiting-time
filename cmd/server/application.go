package main

import (
	"os"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"pizza-waiting-time/config"
)

type application struct {
	config config.Config
	logger *log.Logger
	server *echo.Echo
}

func (app *application) setup() error {
	app.setupLog()
	app.setupConfig()
	app.setupServer()

	return nil
}

func (app *application) setupConfig() {
	config, err := config.Config{}.Load("env")
	if err != nil {
		app.logger.Errorf("error loading config: %v", err)
	}
	app.config = config

	app.logger.Infof("config: %v", config)
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
