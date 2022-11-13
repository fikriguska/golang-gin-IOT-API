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
	statement := "insert into channel (time, value, id_sensor) values (($1), $2, $3)"
	_, err := db.Exec(statement, c.Time, c.Value, c.Id_sensor)
	e.PanicIfNeeded(err)
}
