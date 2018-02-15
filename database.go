package main

import (
	"fmt"
	"log"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

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
	validations := session.DB("guardian").C("validations")

	err := users.DropAllIndexes()
	if err != nil {
		log.Fatal(err)
	}

	index := mgo.Index{
		Key:         []string{"tokenHardExpire"},
		Unique:      true,
		DropDups:    true,
		Background:  true,
		ExpireAfter: 1,
	}
	err = validations.EnsureIndex(index)
	if err != nil {
		fmt.Println(err)
	}

	index = mgo.Index{
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
