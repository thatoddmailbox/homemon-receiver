package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	// Port sets the port to listen for UDP packets.
	Port int
}

var currentConfig config

func loadConfig() error {
	_, err := toml.DecodeFile("config.toml", &currentConfig)
	if os.IsNotExist(err) {
		log.Fatalln("Could not find config file.")
	}
	return err
}
