package auth

import (
	"fmt"
	"net/http"
	"os"
	"wa-service/repo/auth"
	"wa-service/service/wa"

	"github.com/gin-gonic/gin"
)

func LogOut() func(c *gin.Context) {
	return func(c *gin.Context) {
		w := wa.NewWa()
		wac, err := w.Conn()
		nAuth := auth.NewAuth(wac)
		err = nAuth.LogOut()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error logout in: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Error Process",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Success Logout",
		})
	}
}
