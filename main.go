package main

import (
	"fmt"
	"net/http"

	"github.com/kochie/guardian-server/lib"
	//	"./handlers/github"
	//	"./handlers/user"
)

func main() {
	config := lib.ImportConfig()
	addRoutes()

	fmt.Println(config.Mongo.Hostname)
	session := newDatabaseConnection(config.Mongo.Hostname)
	setUpDatabase(session)
	BuildQL(session, config)

	fmt.Println("Server listening on port " + config.Port)
	http.ListenAndServe(":"+config.Port, nil)
	defer session.Close()
}

func addRoutes() {
	//	http.HandleFunc("/api/login/", user.Login)
	//	http.HandleFunc("/api/logout/", user.Logout)

	//	http.HandleFunc("/api/github/connect/", github.Connect)
	//	http.HandleFunc("/api/github/disconnect/", github.Disconnect)
	//	http.HandleFunc("/api/github/addkey/", github.Addkey)
}
