package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
)

func validateHMAC(message, messageMAC, key []byte) bool {
	// from https://golang.org/pkg/crypto/hmac
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func handlePacket(client *net.UDPAddr, data []byte) {
	log.Println(client)
	fmt.Println(hex.Dump(data))
	// 1 byte for battery capacity and power state
	// 2 bytes for battery voltage and power state error flag
	// 8 bytes for timestamp
	// 32 bytes for HMAC-SHA256
	const messageLength = 1 + 2 + 8
	const macLength = 32
	const totalLength = messageLength + macLength
	if len(data) != totalLength {
		// ignore
		return
	}

	// first, validate mac
	message := data[:messageLength]
	mac := data[:macLength]
	if !validateHMAC(message, mac, currentConfig.tokenBytes) {
		// ignore
		return
	}

	log.Println(client)
	fmt.Println(hex.Dump(data))
}

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

	l, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: currentConfig.Port,
	})
	if err != nil {
		panic(err)
	}

	for {
		buf := make([]byte, 1024)
		len, client, err := l.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		go handlePacket(client, buf[0:len])
	}
}
