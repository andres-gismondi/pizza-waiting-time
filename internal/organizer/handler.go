package organizer

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Logger            *log.Logger
	PizzaTimeProvider interface {
		Get(ctx context.Context, pizza int) (time.Duration, error)
	}

	first       *Order
	orders      []*Order
	reorderChan chan []*Order
}

func NewHandler(log *log.Logger) *Handler {
	h := &Handler{
		Logger:      log,
		first:       &Order{Id: -1, Time: time.Second * 100},
		orders:      make([]*Order, 0),
		reorderChan: make(chan []*Order),
	}
	go h.selection()
	return h
}

func (h *Handler) Handle(c echo.Context) error {
	ctx := c.Request().Context()

	var client Client
	if err := json.NewDecoder(c.Request().Body).Decode(&client); err != nil {
		h.Logger.Errorf("could not decode body: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	t, err := h.PizzaTimeProvider.Get(ctx, client.Pizza)
	if err != nil {
		h.Logger.Errorf("could not get pizza: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	go h.organize(ctx, client.Id, t)

	return nil
}

//First approach: log result and endless select
func (h *Handler) organize(ctx context.Context, id int, pizzaTime time.Duration) {
	order := Order{Id: id, Time: pizzaTime}

	if h.first.Id == -1 {
		h.first = &order
	} else {
		h.orders = append(h.orders, &order)
		go h.order()
	}
}

func (h *Handler) selection() {
	defer func() {
		close(h.reorderChan)
	}()

	for {
		select {
		case newOrder := <-h.reorderChan:
			h.Logger.Infof("n: %+v", newOrder)
		case <-time.After(h.first.Time):
			h.Logger.Infof("order [%v] executed", h.first.Id)
			if len(h.orders) != 0 {
				h.first = h.orders[0]
				h.orders = h.orders[1:]
			} else {
				h.first = &Order{Id: -1, Time: 100 * time.Second}
			}
		case <-time.After(time.Second * 100):
			newOrders := make([]*Order, 0)
			for _, order := range h.orders {
				if order.Time <= time.Second {
					order.Time = 0
				} else {
					order.Time = order.Time - time.Second
				}
				newOrders = append(newOrders, order)
			}
			h.orders = newOrders
		}
	}
}

func (h *Handler) order() {
	sort.Slice(h.orders, func(i, j int) bool {
		return h.orders[i].Time.Seconds() < h.orders[j].Time.Seconds()
	})
	h.reorderChan <- h.orders
}
