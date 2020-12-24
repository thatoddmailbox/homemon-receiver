package main

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	// Port sets the port to listen for UDP packets.
	Port int

	// Token sets the token used to authenticate reports.
	Token      string
	tokenBytes []byte
}

var currentConfig config

func loadConfig() error {
	_, err := toml.DecodeFile("config.toml", &currentConfig)
	if os.IsNotExist(err) {
		log.Fatalln("Could not find config file.")
	} else if err != nil {
		return err
	}

	currentConfig.tokenBytes, err = base64.URLEncoding.DecodeString(currentConfig.Token)
	if err != nil {
		return err
	}

	if len(currentConfig.tokenBytes) != 64 {
		log.Fatalf("Token must be exactly 64 bytes long.")
	}

	return nil
}
