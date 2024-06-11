package event

import (
	"encoding/json"
	"go20240606/pkg/mongo"
	"go20240606/pkg/mq"
	"go20240606/pkg/pg"
	"time"

	"github.com/labstack/echo/v4"
)

var MQ mq.IMQ

func init() {
	var err error
	MQ, err = mq.NewIMQ()
	if err != nil {
		panic("mq init error")
	}
}

func RecordEvent(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		res := next(c)
		if res == nil && c.Get("log").(bool) {
			// record the event
			url := c.Get("url").(mongo.Url)
			event := pg.Event{
				Password: url.Password,
				Url:      url.LongURL,
				Time:     time.Now().Format("2006-01-02 15:04:05"),
			}
			// struct to json
			data, err := json.Marshal(event)
			if err != nil {
				println("record event error", err.Error())
			}
			err = MQ.Publish(data)
			if err != nil {
				println("record event error", err.Error())
			}
		}
		return nil
	}
}
