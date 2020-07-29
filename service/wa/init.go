package wa

import (
    "fmt"
    "github.com/Rhymen/go-whatsapp"
    "os"
    "time"
    "wa-service/app"
)

type Wa struct{
    Client *whatsapp.Conn
}

func (w *Wa) Conn() (*whatsapp.Conn, error){
    wac, err := whatsapp.NewConn(5 * time.Second)
    _ = wac.SetClientName(os.Getenv("LongClientName"), os.Getenv("ShortClientName"))
    if err != nil {
        fmt.Println("error creating connection",err)
        return nil,err
    }
    w.Client = wac
    wac.SetClientVersion(app.StoI(os.Getenv("Major")), app.StoI(os.Getenv("Minor")), app.StoI(os.Getenv("Path")))
    return wac,nil
}

func NewWa() *Wa {
    return &Wa{}
}
