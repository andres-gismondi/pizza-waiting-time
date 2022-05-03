package organizer

import "time"

type Client struct {
	Id    int `json:"id"`
	Pizza int `json:"pizza"`
}

type Order struct {
	Id   int
	Time time.Duration
}
