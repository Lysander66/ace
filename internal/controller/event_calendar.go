package controller

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var once sync.Once

type Event struct {
	Id               string `json:"id"`
	Title            string `json:"title"`
	Start            string `json:"start"`
	End              string `json:"end"`
	AllDay           bool   `json:"allDay,omitempty"`
	Editable         bool   `json:"editable,omitempty"`
	StartEditable    bool   `json:"startEditable,omitempty"`
	DurationEditable bool   `json:"durationEditable,omitempty"`
	BackgroundColor  string `json:"backgroundColor,omitempty"`
	TextColor        string `json:"textColor,omitempty"`
}

func Events(c *gin.Context) {
	once.Do(func() {
		for i, event := range events {
			event.Id = strconv.Itoa(i + 1)
			t, err := time.Parse("2006-01-02 15:04", event.Start)
			if err == nil {
				event.End = t.Add(time.Hour).Format("2006-01-02 15:04")
			}
		}
	})
	c.JSON(http.StatusOK, events)
}

var events = []*Event{
	{
		Title: "Reading", BackgroundColor: "green",
		Start: "2023-10-11 21:56:07", End: "2023-10-11 22:56:07",
	},
	{
		Title: "Running",
		Start: "2023-10-12 11:00", End: "2023-10-12 12:00",
	},
	{
		Title: "仙逆",
		Start: "2024-02-05 10:00",
	},
	{
		Title: "炼气十万年",
		Start: "2024-02-06 10:00",
	},
	{
		Title: "少年歌行",
		Start: "2024-02-07 10:00",
	},
	{
		Title: "师兄啊师兄",
		Start: "2024-02-08 10:00",
	},
	{
		Title: "完美世界",
		Start: "2024-02-09 10:00",
	},
	{
		Title: "大主宰年番",
		Start: "2024-02-09 10:00",
	},
	{
		Title: "百炼成神",
		Start: "2024-02-09 10:00",
	},
	{
		Title: "炼气十万年",
		Start: "2024-02-10 10:00",
	},
	{
		Title: "凡人修仙传",
		Start: "2024-02-10 11:00",
	},
	{
		Title: "逆天邪神",
		Start: "2024-02-10 10:00",
	},
	{
		Title: "斗破苍穹年番",
		Start: "2024-02-11 10:00",
	},
	{
		Title: "仙武传",
		Start: "2024-02-11 10:00",
	},
}
