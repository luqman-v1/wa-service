package main

import (
    "github.com/joho/godotenv"
    "log"
    "wa-service/gate"
)



func main() {
    //read file .env
    err := godotenv.Load()
    if err != nil {
        log.Println("Error loading .env file")
    }

    gate.Route()

}


