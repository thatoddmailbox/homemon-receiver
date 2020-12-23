package main

import (
	"log"
)

func main() {
	log.Println("homemon-receiver")

	err := loadConfig()
	if err != nil {
		panic(err)
	}
}
