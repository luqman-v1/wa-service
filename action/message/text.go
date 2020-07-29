package message

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
    "os"
    "wa-service/repo/auth"
    "wa-service/repo/message/text"
    "wa-service/service/wa"
)
type MessageBind struct {
    To string
    Text string
}

func SendMessage() func(c *gin.Context) {
    return func(c *gin.Context) {
        var r MessageBind
        _ = c.ShouldBindJSON(&r);
        w := wa.NewWa()
        wac ,err := w.Conn()
        if err != nil {
            log.Panic(err)
        }
        nAuth := auth.NewAuth(wac)
        _ , err = nAuth.Login()
        if err != nil {
           fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
           c.JSON(http.StatusBadRequest, gin.H{
               "message": "Session Expired",
           })
        }
        t := text.Text{
           To: r.To,
           Message: r.Text,
           Wac:wac,
        }
        go func() {
            t.SendMessage()
        }()
        c.JSON(200, gin.H{
            "message": "Message Sended",
        })
    }
}
