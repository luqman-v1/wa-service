package main

import (
	"log"
	"wa-service/gate"
	"wa-service/service/aws"

	"github.com/joho/godotenv"
)

func main() {
	//read file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	aws.ConnectAws()
	gate.Route()
}
