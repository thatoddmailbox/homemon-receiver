package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
)

func main() {
	generateToken := flag.Bool("generate-token", false, "Generate a random token and exit.")
	flag.Parse()
	if *generateToken {
		newToken := make([]byte, 64)
		_, err := rand.Read(newToken)
		if err != nil {
			panic(err)
		}
		fmt.Println(base64.URLEncoding.EncodeToString(newToken))
		return
	}

	log.Println("homemon-receiver")

	err := loadConfig()
	if err != nil {
		panic(err)
	}
}
