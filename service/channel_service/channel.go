package channel_service

import (
	"src/models"
)

type Channel struct {
	models.Channel
}

func (c *Channel) Add() {
	models.AddChannel(c.Channel)
}
