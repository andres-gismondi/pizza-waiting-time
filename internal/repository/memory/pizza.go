package memory

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type PizzaTime struct {
	Logger *log.Logger
	Order  map[int]time.Duration
}

var (
	ErrNoPizzaTime  = errors.New("pizza time does not exist")
	ErrAddPizzaTime = errors.New("could not add new pizza")
)

func NewPizzaTime(log *log.Logger) *PizzaTime {
	times := make(map[int]time.Duration)
	times[1] = time.Second * 2
	times[2] = time.Second * 3
	times[3] = time.Second * 5

	return &PizzaTime{Logger: log, Order: times}
}

func (pt *PizzaTime) Get(ctx context.Context, pizza int) (time.Duration, error) {
	//pt.Logger.Infof("get pizza [%v] time", pizza)

	time, ok := pt.Order[pizza]
	if !ok {
		return 0, fmt.Errorf("pizza time does not exist with pizza: %v: %w", pizza, ErrNoPizzaTime)
	}

	return time, nil
}

func (pt *PizzaTime) Add(ctx context.Context, pizza int, duration time.Duration) error {
	//pt.Logger.Infof("adding pizza [%v] with time [%v]", pizza, duration)

	pt.Order[pizza] = duration
	if _, ok := pt.Order[pizza]; !ok {
		return fmt.Errorf("could not add new pizza [%v] with time [%v]: %w", pizza, duration, ErrAddPizzaTime)
	}

	return nil
}
