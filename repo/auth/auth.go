package auth

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
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
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			session, err = a.wac.Login(qr)
			if err != nil {
				_ = fmt.Errorf("error during login: %v\n", err)
				return
			}
			//save session
			err = a.writeSession(session)
			if err != nil {
				_ = fmt.Errorf("error saving session: %v\n", err)
				return
			}
			return
		}()
		_ = qrcode.WriteFile(<-qr, qrcode.Medium, 256, app.PATH_FILE+filename)
	}

	return filename, nil
}

func (a *Auth) readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	var file io.Reader
	if os.Getenv("FILE_MANAGER") != "s3" {
		file, err := os.Open(app.PATH_SESSION)
		if err != nil {
			return session, err
		}
		defer file.Close()
	} else {
		b, err := s3.Get(app.FILENAME_SESSION)
		if err != nil {
			return session, err
		}
		file = bytes.NewReader(b)
	}
	decoder := gob.NewDecoder(file)
	err := decoder.Decode(&session)
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
	err := s3.Delete(app.FILENAME_SESSION)
	if err != nil {
		log.Println("Failed Delete Session on S3", err)
		return err
	}
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
