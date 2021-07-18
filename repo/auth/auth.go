package auth

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"wa-service/app"
	"wa-service/entity"

	"github.com/Rhymen/go-whatsapp"
	"github.com/skip2/go-qrcode"
)

const QRCODE_FILENAME = "qr.png"

type Auth struct {
	wac *whatsapp.Conn
}

func NewAuth(wac *whatsapp.Conn) *Auth {
	return &Auth{
		wac: wac,
	}
}

//CheckSession wa
func (a *Auth) CheckSession() error {
	session, err := a.readSession()
	if err != nil {
		log.Println("Failed restore session", err)
		return err
	}
	//restore session
	session, err = a.wac.RestoreWithSession(session)
	if err != nil {
		err = a.deleteSession()
		if err != nil {
			log.Println("Failed Delete session", err)
			return err
		}
	}
	return nil
}

//Login and show qr code
func (a *Auth) Login() (entity.Auth, error) {
	filename := QRCODE_FILENAME
	//load saved session
	session, err := a.readSession()
	if err == nil {
		//restore session
		session, err = a.wac.RestoreWithSession(session)
		if err != nil {
			err = a.deleteSession()
			if err != nil {
				log.Println("Failed Delete session", err)
				return entity.Auth{}, err
			}
			return a.Login()
		}
		return entity.Auth{}, app.ErrorAlreadyScanned()
	} else {
		log.Println("create qr code and session ...")
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			err := func() error {
				session, err = a.wac.Login(qr)
				log.Println("login ...", session, err)
				if err != nil {
					_ = fmt.Errorf("error during login: %v\n", err)
					return err
				}
				//save session
				err = a.writeSession(session)
				if err != nil {
					_ = fmt.Errorf("error saving session: %v\n", err)
					return err
				}
				return err
			}()
			if err != nil {
				return
			}
		}()
		if err = qrcode.WriteFile(<-qr, qrcode.Medium, 256, app.PATH_FILE+filename); err != nil {
			log.Println("err at write file ...", err)
			return entity.Auth{}, err
		}
	}
	return entity.Auth{FileName: filename}, nil
}

//read session from local file
func (a *Auth) readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(app.PATH_SESSION)
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	if err = decoder.Decode(&session); err != nil {
		return session, err
	}
	return session, nil
}

//write session
func (a *Auth) writeSession(session whatsapp.Session) error {
	file, err := os.Create(app.PATH_SESSION)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	if err = encoder.Encode(session); err != nil {
		return err
	}
	return nil
}

//delete session file gob on local storage
func (a *Auth) deleteSession() error {
	return os.Remove(app.PATH_SESSION)
}

// LogOut logout session wa
func (a *Auth) LogOut() error {
	if err := a.wac.Logout(); err != nil {
		return err
	}
	if err := a.deleteSession(); err != nil {
		return err
	}
	return nil
}
