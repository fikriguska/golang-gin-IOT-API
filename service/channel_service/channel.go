package channel_service

import (
	"src/models"
	"time"
)

type Channel struct {
	models.Channel
}

func (c *Channel) Add() {
	t := time.Now()
	c.Time = t
	models.AddChannel(c.Channel)
}
