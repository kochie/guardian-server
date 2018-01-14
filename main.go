package main

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	//	"./handlers/github"
	//	"./handlers/user"
)

func main() {
	port := ":8080"

	addRoutes()

	//	github.CreateConnection()
	setup()
	fmt.Println("Server listening on port " + port[1:])
	http.ListenAndServe(port, nil)
}

func addRoutes() {
	//	http.HandleFunc("/api/login/", user.Login)
	//	http.HandleFunc("/api/logout/", user.Logout)

	//	http.HandleFunc("/api/github/connect/", github.Connect)
	//	http.HandleFunc("/api/github/disconnect/", github.Disconnect)
	//	http.HandleFunc("/api/github/addkey/", github.Addkey)
}

func setup() {
	session := newDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows)
}

func newDatabaseConnection() (session *mgo.Session) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	return
}
