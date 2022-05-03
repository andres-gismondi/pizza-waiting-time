package organizer_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"pizza-waiting-time/internal/organizer"
)

func TestHandler_Handle(t *testing.T) {
	logg := log.New()
	logg.SetFormatter(&log.JSONFormatter{})
	logg.SetLevel(log.DebugLevel)
	logg.SetOutput(os.Stdout)
	repo := organizer.NewPizzaTime(logg)
	handler := organizer.NewHandler(logg)
	handler.PizzaTimeProvider = repo

	server := echo.New()

	request := `{
				"id": 1,
				"pizza": 3
			}`
	req := httptest.NewRequest(http.MethodPost, "/pizza/", strings.NewReader(request))
	ctx := server.NewContext(req, nil)
	handler.Handle(ctx)

	request = `{
				"id": 2,
				"pizza": 3
			}`
	req = httptest.NewRequest(http.MethodPost, "/pizza/", strings.NewReader(request))
	ctx = server.NewContext(req, nil)
	handler.Handle(ctx)

	request = `{
				"id": 3,
				"pizza": 2
			}`
	req = httptest.NewRequest(http.MethodPost, "/pizza/", strings.NewReader(request))
	ctx = server.NewContext(req, nil)
	time.Sleep(2 * time.Second)
	handler.Handle(ctx)

	request = `{
				"id": 4,
				"pizza": 1
			}`
	req = httptest.NewRequest(http.MethodPost, "/pizza/", strings.NewReader(request))
	ctx = server.NewContext(req, nil)
	time.Sleep(2 * time.Second)
	handler.Handle(ctx)

	time.Sleep(19 * time.Second)
}
