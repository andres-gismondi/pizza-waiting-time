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
	NotificationProvider interface {
		Notify(order Order)
	}

	first       *Order
	orders      []*Order
	reorderChan chan []*Order
}

func NewHandler(log *log.Logger) *Handler {
	h := &Handler{
		Logger: log,
		first:  &Order{Id: -1, Time: time.Second * 5},
		orders: make([]*Order, 0),
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
		h.Logger.Infof("first order [%v]", h.first.Id)
	} else {
		h.orders = append(h.orders, &order)
		go h.order()
	}
}

func (h *Handler) selection() {
	for {
		select {
		case <-time.After(h.first.Time):
			if len(h.orders) != 0 || h.first.Id != -1 {
				h.Logger.Infof("order [%v] executed", h.first.Id)

				h.first = h.orders[0]
				h.orders = h.orders[1:]

				h.removeOrder(h.first.Id)
			} else {
				h.first = &Order{Id: -1, Time: time.Second * 10}
			}
		case <-time.After(time.Second * 3):
			if len(h.orders) != 0 || h.first.Id != -1 {
				newOrders := make([]*Order, 0)
				for _, order := range h.orders {
					if order.Time <= time.Second {
						h.removeOrder(order.Id)
					} else {
						order.Time = order.Time - time.Second
					}
					newOrders = append(newOrders, order)
				}
				h.orders = newOrders
				h.Logger.Infof("times: %v", h.orders)
				if len(h.orders) == 0 {
					h.first = &Order{Id: -1, Time: time.Second * 10}
				}
			} else {
				h.first = &Order{Id: -1, Time: time.Second * 10}
			}
		}
	}
}

func (h *Handler) order() {
	sort.Slice(h.orders, func(i, j int) bool {
		return h.orders[i].Time.Seconds() < h.orders[j].Time.Seconds()
	})
	h.Logger.Infof("after: %v", h.orders)
}

func (h *Handler) removeOrder(id int) {
	index := -1
	for i, or := range h.orders {
		if id == or.Id {
			index = i
		}
	}
	if index == -1 {
		return
	}

	h.NotificationProvider.Notify(*h.first)
	h.Logger.Infof("order with id:[%v] removed", id)
	h.orders = append(h.orders[:index], h.orders[index+1:]...)
}
