package auth

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "net/http"
    "os"
    "wa-service/app"
    "wa-service/repo/auth"
    "wa-service/service/wa"
)

func Login() func(c *gin.Context) {
    return func(c *gin.Context) {
        w := wa.NewWa()
        wac ,err := w.Conn()
        nAuth := auth.NewAuth(wac)
        filename , err := nAuth.Login()
        if err != nil {
            fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
            return
        }
        c.JSON(200, gin.H{
            "message": http.Dir(app.PATH_FILE + filename),
        })
    }
}
