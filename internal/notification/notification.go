package notification

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"pizza-waiting-time/internal/organizer"
)

type Notification struct {
	Logger     *log.Logger
	HTTPPoster interface {
		Post(ctx context.Context, path string, body interface{}) (*http.Response, error)
	}
}

func (nf *Notification) Notify(order organizer.Order) {
	nf.Logger.Infof("notifing order: %v", order)

	path := "/notification"
	_, err := nf.HTTPPoster.Post(context.Background(), path, order)
	if err != nil {
		nf.Logger.Errorf("could not send notification")
	}
}
