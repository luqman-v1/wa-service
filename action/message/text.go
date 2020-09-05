package message

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"wa-service/repo/auth"
	"wa-service/repo/message/text"
	"wa-service/service/wa"

	"github.com/gin-gonic/gin"
)

type MessageBind struct {
	To   string
	Text string
}

func SendMessage() func(c *gin.Context) {
	return func(c *gin.Context) {
		var r MessageBind
		err := c.ShouldBindJSON(&r)
		if err != nil {
			log.Println(err)
		}
		w := wa.NewWa()
		wac, _ := w.Conn()
		nAuth := auth.NewAuth(wac)
		err = nAuth.CheckSession()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Session Expired",
			})
			return
		}
		log.Println("request ==>", r)
		t := text.Text{
			To:      r.To,
			Message: r.Text,
			Wac:     wac,
		}
		go func() {
			err := t.SendMessage()
			if err != nil {
				log.Println(err)
			}
		}()

		c.JSON(http.StatusOK, gin.H{
			"message": "Message has been processed",
		})
		return
	}
}
