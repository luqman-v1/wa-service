package app

import "errors"

func ErrorAlreadyScanned() error {
	return errors.New("QR Code Has Been Scan")
}
