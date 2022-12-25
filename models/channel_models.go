package models

import (
	e "src/error"
	"time"
)

type Channel struct {
	Id        int
	Time      time.Time
	Value     float64
	Id_sensor int
}

type ChannelAdd struct {
	Value     float64 `json:"value" binding:"required"`
	Id_sensor int     `json:"id_sensor" binding:"required"`
}

func AddChannel(c Channel) {
	statement := replaceQueryParam("insert into channel (time, value, id_sensor) values (now(), %g, %d)", c.Value, c.Id_sensor)
	_, err := db.Exec(statement)
	e.PanicIfNeeded(err)
}
