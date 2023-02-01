package models

import (
	e "src/error"
	"time"
)

type Channel struct {
	Id      int
	Time    time.Time
	Value   []*float64
	Id_node int
}

type ChannelGet struct {
	Value []*float64 `json:"value"`
	Time  time.Time  `json:"time"`
}

type ChannelAdd struct {
	Value   []*float64 `json:"value" binding:"required"`
	Id_node int        `json:"id_node" binding:"required"`
}

func AddChannel(c Channel) {
	statement := "insert into feed (time, value, id_node) values (($1), $2, $3)"
	_, err := db.Exec(cb(), statement, c.Time, c.Value, c.Id_node)
	e.PanicIfNeeded(err)
}

func GetFeedByNodeId(id int, limit int) []ChannelGet {
	var feed ChannelGet
	var feeds []ChannelGet
	feeds = make([]ChannelGet, 0)
	statement := "select time, value from feed where id_node = $1 ORDER BY time DESC limit $2;"
	rows, err := db.Query(cb(), statement, id, limit)
	e.PanicIfNeeded(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&feed.Time, &feed.Value)
		e.PanicIfNeeded(err)
		feeds = append(feeds, feed)
	}
	return feeds
}
