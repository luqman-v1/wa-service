package wa

import (
	"log"
	"os"
	"time"
	"wa-service/app"

	"github.com/Rhymen/go-whatsapp"
)

type Wa struct {
	Client *whatsapp.Conn
}

func (w *Wa) Conn() (*whatsapp.Conn, error) {
	version := os.Getenv("Major") + "." + os.Getenv("Minor") + "." + os.Getenv("Path")
	wac, err := whatsapp.NewConn(5 * time.Second)
	err = wac.SetClientName(os.Getenv("LongClientName"), os.Getenv("ShortClientName"), version)
	if err != nil {
		log.Fatal("error creating connection", err)
		return nil, err
	}
	w.Client = wac
	wac.SetClientVersion(app.StoI(os.Getenv("Major")), app.StoI(os.Getenv("Minor")), app.StoI(os.Getenv("Path")))
	return wac, nil
}

func NewWa() *Wa {
	return &Wa{}
}
