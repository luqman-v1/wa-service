package auth

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"wa-service/app"
	"wa-service/repo/s3"

	"github.com/Rhymen/go-whatsapp"
	"github.com/skip2/go-qrcode"
)

type Auth struct {
	wac *whatsapp.Conn
}

func NewAuth(wac *whatsapp.Conn) *Auth {
	return &Auth{
		wac: wac,
	}
}

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

func (a *Auth) Login() (string, error) {
	filename := "qr.png"
	//load saved session
	session, err := a.readSession()
	if err == nil {
		//restore session
		session, err = a.wac.RestoreWithSession(session)
		if err != nil {
			err = a.deleteSession()
			if err != nil {
				log.Println("Failed Delete session", err)
				return "", err
			}
			return a.Login()
		}
		return "", errors.New("QR Code Has Been Scan")
	} else {
		log.Println("create qc code and session ...")
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
		err = qrcode.WriteFile(<-qr, qrcode.Medium, 256, app.PATH_FILE+filename)
		if err != nil {
			log.Println("err at writefile ...", err)
			return "", err
		}
	}

	return filename, nil
}

func (a *Auth) readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(app.PATH_SESSION)
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
	_, _ = s3.Upload(app.FILENAME_SESSION, app.PATH_SESSION)
	return nil
}

func (a *Auth) deleteSession() error {
	return os.Remove(app.PATH_SESSION)
}

func (a *Auth) LogOut() error {
	err := a.wac.Logout()
	if err != nil {
		return err
	}
	_ = a.deleteSession()
	return nil
}
