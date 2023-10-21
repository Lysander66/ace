package controller

import (
	"github.com/Lysander66/ace/internal/service"
	"github.com/gin-gonic/gin"
)

func SetRouter(g *gin.Engine) {
	g.GET("/events", Events)

	api := g.Group("/api/v1")
	api.GET("crypto/md5", md5Sum)
	api.GET("crypto/sha1", sha1Sum)
	api.GET("crypto/sha256", sha256Sum)
	api.GET("crypto/sha512", sha512Sum)

	// WebSocket
	wsRouter := g.Group("/ws/v1")
	wsRouter.GET("echo", service.Echo)
}
