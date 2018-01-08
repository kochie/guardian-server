package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"

	"./handlers/github"
	"./handlers/user"
)

func main() {
	port := ":8080"

	addRoutes()

	github.CreateConnection()

	fmt.Println("Server listening on port " + port[1:])
	http.ListenAndServe(port, nil)
}

func addRoutes() {
	http.HandleFunc("/api/login/", user.Login)
	http.HandleFunc("/api/logout/", user.Logout)

	http.HandleFunc("/api/github/connect/", github.Connect)
	http.HandleFunc("/api/github/disconnect/", github.Disconnect)
	http.HandleFunc("/api/github/addkey/", github.Addkey)
}
