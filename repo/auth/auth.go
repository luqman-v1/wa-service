package auth

import (
    "encoding/gob"
    "fmt"
    "github.com/Rhymen/go-whatsapp"
    "github.com/skip2/go-qrcode"
    "os"
    "wa-service/app"
)

type Auth struct{
    wac *whatsapp.Conn
}

func NewAuth(wac *whatsapp.Conn) *Auth {
    return &Auth{
        wac :wac,
    }
}

func (a *Auth) Login() (string,error) {
    filename := "qr.png"
    //load saved session
    session, err := a.readSession()
    if err == nil {
        //restore session
        session, err = a.wac.RestoreWithSession(session)
        if err != nil {
            err = a.deleteSession()
            if err != nil {
               return "",err
            }
            return a.Login()
        }
    } else {
        //no saved session -> regular login
        qr := make(chan string)
        go func() {
            session, err = a.wac.Login(qr)
            if err != nil {
               fmt.Errorf("error during login: %v\n", err)
                return
            }
            //save session
            err = a.writeSession(session)
            if err != nil {
                 fmt.Errorf("error saving session: %v\n", err)
                 return
            }
            return
        }()
        _ = qrcode.WriteFile(<-qr, qrcode.Medium, 256, app.PATH_FILE + filename)
    }


    return filename , nil
}

func (a *Auth) readSession() (whatsapp.Session, error) {
    session := whatsapp.Session{}
    file, err := os.Open(  app.PATH_SESSION)
    if err != nil {
        return session, err
    }
    defer file.Close()
    decoder := gob.NewDecoder(file)
    err = decoder.Decode(&session)
    if err != nil {
        return session, err
    }
    return session, nil
}

func (a *Auth) writeSession(session whatsapp.Session) error {
    file, err := os.Create(app.PATH_SESSION)
    if err != nil {
        return err
    }
    defer file.Close()
    encoder := gob.NewEncoder(file)
    err = encoder.Encode(session)
    if err != nil {
        return err
    }
    return nil
}

func (a *Auth) deleteSession() error {
    return os.Remove(app.PATH_SESSION)
}

func (a *Auth) LogOut () error{
    err := a.wac.Logout()
    if err != nil {
        return err
    }
    _ = a.deleteSession()
    return nil
}