package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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
	BuildQL(session)

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

func newDatabaseConnection(hostname string) (session *mgo.Session) {
	session, err := mgo.Dial(hostname)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func setUpDatabase(session *mgo.Session) {
	session.SetMode(mgo.Monotonic, true)
	users := session.DB("guardian").C("users")

	err := users.DropAllIndexes()
	if err != nil {
		log.Fatal(err)
	}

	index := mgo.Index{
		Key:           []string{"number"},
		Unique:        true,
		DropDups:      true,
		Background:    true,
		PartialFilter: bson.M{"number": bson.M{"$type": "string"}},
	}
	err = users.EnsureIndex(index)
	if err != nil {
		fmt.Println(err)
	}

	index = mgo.Index{
		Key:           []string{"email"},
		Unique:        true,
		DropDups:      true,
		Background:    true,
		PartialFilter: bson.M{"email": bson.M{"$type": "string"}},
	}
	err = users.EnsureIndex(index)
	if err != nil {
		fmt.Println(err)
	}

	index = mgo.Index{
		Key:           []string{"username"},
		Unique:        true,
		DropDups:      true,
		Background:    true,
		PartialFilter: bson.M{"username": bson.M{"$type": "string"}},
	}
	err = users.EnsureIndex(index)
	if err != nil {
		fmt.Println(err)
	}

	index = mgo.Index{
		Key:        []string{"services.name", "number", "email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = users.EnsureIndex(index)
	if err != nil {
		fmt.Println(err)
	}

	index = mgo.Index{
		Key:        []string{"devices.token", "devices.name", "number", "email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = users.EnsureIndex(index)
	if err != nil {
		fmt.Println(err)
	}
}
