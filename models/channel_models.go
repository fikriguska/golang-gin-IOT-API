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
	statement := "insert into channel (time, value, id_sensor) values (now(), $1, $2)"
	_, err := db.Exec(cb(), statement, c.Value, c.Id_sensor)
	e.PanicIfNeeded(err)
}
