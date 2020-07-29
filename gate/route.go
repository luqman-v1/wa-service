package gate
import (
    "github.com/gin-gonic/gin"
    "wa-service/action/auth"
    "wa-service/action/message"
    "wa-service/app"
)

func Route() {
    r := gin.New()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Use(Middleware())
    r.Static("/storage/file", app.PATH_FILE)
    r.GET("/login", auth.Login())
    r.POST("/send/message", message.SendMessage())
    r.Use(gin.Recovery())
    _ = r.Run()
}