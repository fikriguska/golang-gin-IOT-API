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
	// timeStamp := t.Format("15:04:05 UTC")

	c.Time = t
	models.AddChannel(c.Channel)
}
