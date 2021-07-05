package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"wa-service/app"
	"wa-service/repo/auth"
	"wa-service/service/wa"

	"github.com/gin-gonic/gin"
)

func Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		w := wa.NewWa()
		wac, err := w.Conn()
		nAuth := auth.NewAuth(wac)
		filename, err := nAuth.Login()
		log.Println("err filename", err)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": http.Dir(app.PATH_FILE + filename),
		})
	}
}
