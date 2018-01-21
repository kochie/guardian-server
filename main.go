package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/kochie/guardian-server/lib"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//	"./handlers/github"
	//	"./handlers/user"
)

//User datatype
type User struct {
	Email  string `json:"email"`
	Number string `json:"number"`
}

type Service struct {
	Name   string
	Secret string
	UserID string
	Active bool
}

type Device struct {
	Name   string
	UserID string
	Active bool
	Token  string
}

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"number": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func searchForUser(session *mgo.Session, user *User) (*User, err error) {
	users := session.DB("guardian").C("users")
	err = users.Find(bson.M{"number": user.Number})
	return &User, err
}

func buildQL(session *mgo.Session) {
	users := session.DB("guardian").C("users")

	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"login": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				login := p.Args["login"].(string)

				result := User{}
				err := users.Find(bson.M{"email": login}).One(&result)
				if err != nil {
					// log.Fatalf("Failed to read from the database: %v", err)
					log.Println("No email found looking for number")
					err = users.Find(bson.M{"number": login}).One(&result)
					if err != nil {
						log.Println("No number found! User does not exist")
						return nil, errors.New("No user found")
					}
				}

				return result, nil
			},
		},
	}
	mutations := graphql.Fields{
		"addUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"number": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var email, number string

				if p.Args["number"] == nil && p.Args["email"] == nil {
					return nil, errors.New("Either an email or phone number is required")
				}

				if p.Args["email"] != nil {
					email = p.Args["email"].(string)
				}
				if p.Args["number"] != nil {
					number = p.Args["number"].(string)
				}

				user := &User{Email: email, Number: number}
				searchForUser(user)
				err := users.Insert(user)
				if err != nil {
					return nil, err
				}
				return user, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutations}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery), Mutation: graphql.NewObject(rootMutation)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)

}

func main() {
	config := lib.ImportConfig()
	addRoutes()

	session := setup(config.MongoHostname)
	buildQL(session)

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

func setup(hostname string) (session *mgo.Session) {
	session = newDatabaseConnection(hostname)
	session.SetMode(mgo.Monotonic, true)

	// users := session.DB("guardian").C("users")
	// err := users.Insert(&User{Email: "robert@kochie.io"})

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// result := User{}
	// err = users.Find(bson.M{"email": "robert@kochie.io"}).One(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("email:", result.Email)

	return session
}

func newDatabaseConnection(hostname string) (session *mgo.Session) {
	session, err := mgo.Dial(hostname)
	if err != nil {
		log.Fatal(err)
	}
	return
}
