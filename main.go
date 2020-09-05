package main

import (
	"log"
	"wa-service/gate"

	"github.com/joho/godotenv"
)

func main() {
	//read file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	gate.Route()
}
