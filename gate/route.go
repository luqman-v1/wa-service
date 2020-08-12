package gate

import (
	"wa-service/action/auth"
	"wa-service/action/message"
	"wa-service/app"

	"github.com/gin-gonic/gin"
)

func Route() {
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Static("/storage/file", app.PATH_FILE)
	r.Use(Middleware())
	r.GET("/login", auth.Login())
	r.POST("/send/message", message.SendMessage())
	r.Use(gin.Recovery())
	_ = r.Run()
}
