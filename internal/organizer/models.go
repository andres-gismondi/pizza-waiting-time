package organizer

import (
	"fmt"
	"time"
)

type Client struct {
	Id    int `json:"id"`
	Pizza int `json:"pizza"`
}

type Order struct {
	Id   int
	Time time.Duration
}

func (o Order) String() string {
	return fmt.Sprintf("{id:%v | time:%v}", o.Id, o.Time)
}
