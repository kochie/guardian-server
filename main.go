package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kochie/guardian-server/lib"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//	"./handlers/github"
	//	"./handlers/user"
)

//User datatype
type User struct {
	Email string
}

func main() {
	config := lib.ImportConfig()
	addRoutes()

	setup(config.MongoHostname)
	fmt.Println("Server listening on port " + config.Port)
	http.ListenAndServe(":"+config.Port, nil)
}

func addRoutes() {
	//	http.HandleFunc("/api/login/", user.Login)
	//	http.HandleFunc("/api/logout/", user.Logout)

	//	http.HandleFunc("/api/github/connect/", github.Connect)
	//	http.HandleFunc("/api/github/disconnect/", github.Disconnect)
	//	http.HandleFunc("/api/github/addkey/", github.Addkey)
}

func setup(hostname string) {
	session := newDatabaseConnection(hostname)
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	users := session.DB("guardian").C("users")
	err := users.Insert(&User{"robert@kochie.io"})

	if err != nil {
		log.Fatal(err)
	}

	result := User{}
	err = users.Find(bson.M{"email": "robert@kochie.io"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("email:", result.Email)
}

func newDatabaseConnection(hostname string) (session *mgo.Session) {
	session, err := mgo.Dial(hostname)
	if err != nil {
		log.Fatal(err)
	}
	return
}
