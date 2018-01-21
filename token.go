package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func GenerateToken() (token string) {
	c := 256
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("error:", err)
	}

	token = base64.StdEncoding.EncodeToString(b)
	return
}
