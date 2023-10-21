package service

import (
	"log/slog"

	"github.com/Lysander66/ace/pkg/ws"
	"github.com/gin-gonic/gin"
)

var hub = ws.NewHub()

func Echo(c *gin.Context) {
	var (
		header = c.Request.Header
		ip     = header.Get("X-Forwarded-For")
	)
	if header.Get("Cdn-Loop") == "cloudflare" {
		ip = header.Get("Cf-Connecting-Ip")
		slog.Debug("ClientIP", "X-Forwarded-For", header.Get("X-Forwarded-For"), "Cf-Connecting-Ip", ip)
	}

	ws.ServeWs(c.Writer, c.Request, hub)
}
