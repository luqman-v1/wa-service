package text

import (
    "fmt"
    "github.com/Rhymen/go-whatsapp"
    "time"
    "wa-service/app"
)

type Text struct {
    To string
    Message string
    Wac *whatsapp.Conn
}

func (t *Text)SendMessage(){
    <-time.After(3 * time.Second)

    msg := whatsapp.TextMessage{
        Info: whatsapp.MessageInfo{
            RemoteJid: app.WidToJid(t.To),
        },
        ContextInfo: whatsapp.ContextInfo{},
        Text:        t.Message,
    }

    msgId, err := t.Wac.Send(msg)
    if err != nil {
        fmt.Println("error sending message : " , err)
    } else {
        fmt.Println("Message Sent -> ID : " + msgId)
    }
}