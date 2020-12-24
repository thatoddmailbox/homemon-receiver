package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

const batteryCapacityError = 0x7F
const batteryVoltageError = 0x1FFF

type usbStatus int8

const (
	usbStatusNotPresent usbStatus = 0
	usbStatusPresent    usbStatus = 1
	usbStatusError      usbStatus = -1
)

func validateHMAC(message, messageMAC, key []byte) bool {
	// from https://golang.org/pkg/crypto/hmac
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func handlePacket(client *net.UDPAddr, data []byte) {
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
	mac := data[messageLength : messageLength+macLength]
	if !validateHMAC(message, mac, currentConfig.tokenBytes) {
		// ignore
		return
	}

	// see homemon-daemon's transportUDP.go for clearer documentation about the format
	batteryCapacity := message[0] & 0x7F
	usbPresent := (message[0] & 0x80) != 0
	batteryVoltage := binary.BigEndian.Uint16(message[1:]) & 0x1FFF
	usbError := (message[1] & 0x80) != 0
	messageTimestamp := int64(binary.BigEndian.Uint64(message[3:]))
	currentTimestamp := time.Now().Unix()

	if currentTimestamp-messageTimestamp > 10 || currentTimestamp-messageTimestamp < -10 {
		// it's a difference of more than 10 seconds, consider the message expired
		log.Printf("Got expired message from %s, ignoring!", client.String())
		return
	}

	var powered usbStatus
	if usbPresent && !usbError {
		powered = usbStatusPresent
	} else if !usbPresent && !usbError {
		powered = usbStatusNotPresent
	} else {
		powered = usbStatusError
	}

	batteryCapacityDatabase := int8(batteryCapacity)
	if batteryCapacity == batteryCapacityError {
		batteryCapacityDatabase = -1
	}

	batteryVoltageDatabase := int16(batteryVoltage)
	if batteryVoltage == batteryVoltageError {
		batteryVoltageDatabase = -1
	}

	_, err := db.Exec(
		"INSERT INTO reports(powered, batteryLevel, batteryVoltage, ip, transport, clientTimestamp, timestamp) VALUES(?, ?, ?, ?, ?, ?, ?)",
		powered,
		batteryCapacityDatabase,
		batteryVoltageDatabase,
		client.IP.String(),
		1,
		messageTimestamp,
		currentTimestamp,
	)
	if err != nil {
		log.Printf("Encountered error when processing message from %s!", client.String())
		log.Println(err)
		return
	}
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

	err = connectToDatabase()
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
