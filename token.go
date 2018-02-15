package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(bytesize int) (token string) {
	b := make([]byte, bytesize)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("error:", err)
	}

	token = base64.StdEncoding.EncodeToString(b)
	return
}

func JWT() {
	// sample token string taken from the New example
	var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	hmacSampleSecret := []byte("A secret so powerful")

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	var token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	// token = jwt.NewWithClaims(jwt.SigningMethodECDSA, jwt.MapClaims{
	// 	"foo": "bar",
	// 	"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	// })

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err = token.SignedString(hmacSampleSecret)

	fmt.Println(tokenString, err)
}

func makeToken() string {
	return "Cool Whip"
}
