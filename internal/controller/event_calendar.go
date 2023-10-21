package controller

import (
	"net/http"
	"strconv"
	"sync"

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
		}
	})
	c.JSON(http.StatusOK, events)
}

var events = []*Event{
	{
		Title:           "Reading",
		Start:           "2023-10-11 21:56:07",
		End:             "2023-10-11 22:56:07",
		BackgroundColor: "green",
	},
	{
		Title: "Running",
		Start: "2023-10-12 11:00",
		End:   "2023-10-12 12:00",
	},
}
