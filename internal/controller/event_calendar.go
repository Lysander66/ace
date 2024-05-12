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
}
