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

// Service is the database structure of service
type Service struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Active      bool   `json:"active"`
	Description string `json:"description"`
}

// Device is the database structure of device
type Device struct {
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	Token       string `json:"token"`
	Description string `json:"description"`
}

//User datatype
type User struct {
	Email    string    `json:"email"`
	Number   string    `json:"number"`
	Services []Service `json:"services"`
	Devices  []Device  `json:"devices"`
}

func getServices(userLogin string, p graphql.ResolveParams, users *mgo.Collection) ([]Service, error) {
	result := User{}
	query := bson.M{"$or": []bson.M{{"email": userLogin}, {"number": userLogin}}}
	response := bson.M{"services": 1}

	if p.Args["filter"] != nil {
		filter := p.Args["filter"].(map[string]interface{})

		if len(filter) > 0 {
			if p.Args["intersection"] == false {
				params := make([]bson.M, 1)
				for e, r := range filter {
					params = append(params, bson.M{e: r})
				}
				query["services"] = bson.M{"$elemMatch": bson.M{"$or": params}}
			} else {
				query["services"] = bson.M{"$elemMatch": filter}
			}
		}
	}

	err := users.Find(query).Select(response).One(&result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result.Services, nil
}

func buildQL(session *mgo.Session) {
	users := session.DB("guardian").C("users")

	var serviceType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Service",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The name given to this particular service",
			},
			"type": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The type of service, specified by the provider.",
			},
			"active": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "The state of this service, whether it is active or not.",
			},
		},
	})

	var deviceType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Device",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"active": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"token": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
	})

	serviceFilter := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "serviceFilter",
			Fields: graphql.InputObjectConfigFieldMap{
				"name": &graphql.InputObjectFieldConfig{
					Type:        graphql.String,
					Description: "The name of a service associated with this user",
				},
				"type": &graphql.InputObjectFieldConfig{
					Type:        graphql.String,
					Description: "The service provider of a service associated with this user",
				},
				"active": &graphql.InputObjectFieldConfig{
					Type:        graphql.Boolean,
					Description: "The state of a service associated with this user",
				},
			},
		},
	)

	deviceInput := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "DeviceInput",
			Fields: graphql.InputObjectConfigFieldMap{
				"name": &graphql.InputObjectFieldConfig{
					Type:        graphql.String,
					Description: "The name of a device associated with this user",
				},
				"token": &graphql.InputObjectFieldConfig{
					Type:        graphql.String,
					Description: "The device token associated with this user",
				},
				"active": &graphql.InputObjectFieldConfig{
					Type:        graphql.Boolean,
					Description: "The state of a service associated with this user",
				},
			},
		},
	)

	// fmt.Println(p.Source.(User).Email)
	var userType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"number": &graphql.Field{
				Type: graphql.String,
			},
			"devices": &graphql.Field{
				Type: &graphql.List{OfType: deviceType},
			},
			"services": &graphql.Field{
				Type: &graphql.List{OfType: serviceType},
				Args: graphql.FieldConfigArgument{
					"filter": &graphql.ArgumentConfig{
						Type: serviceFilter,
					},
					"intersection": &graphql.ArgumentConfig{
						Type:         graphql.Boolean,
						Description:  "If the search arguments should be accepted in union or intersection.",
						DefaultValue: false,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := p.Source.(User)
					if user.Number != "" {
						return getServices(user.Number, p, users)
					}
					return getServices(user.Email, p, users)
				},
			},
		},
	})

	userInputType := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "userInput",
			Fields: graphql.InputObjectConfigFieldMap{
				"email": &graphql.InputObjectFieldConfig{
					Type: graphql.String,
				},
				"number": &graphql.InputObjectFieldConfig{
					Type: graphql.String,
				},
			},
		},
	)

	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: &graphql.List{OfType: graphql.String},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return []string{"hello", "hello"}, nil
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
				query := bson.M{"$or": []bson.M{{"email": login}, {"number": login}}}
				err := users.Find(query).One(&result)
				if err != nil {
					return nil, errors.New("No user found")
				}
				return result, nil
			},
		},
		"services": &graphql.Field{
			Type: &graphql.List{OfType: serviceType},
			Args: graphql.FieldConfigArgument{
				"login": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The login of the user to search services against.",
				},
				"filter": &graphql.ArgumentConfig{
					Type: serviceFilter,
				},
				"intersection": &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					Description:  "If the search arguments should be accepted in union or intersection.",
					DefaultValue: false,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userLogin := p.Args["login"].(string)
				return getServices(userLogin, p, users)
			},
		},
		"devices": &graphql.Field{
			Type: &graphql.List{OfType: deviceType},
			Args: graphql.FieldConfigArgument{
				"login": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The login of the user to search devices against.",
				},
				"name": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "The name of a device associated with this user",
				},
				"token": &graphql.ArgumentConfig{
					Type:        graphql.String,
					Description: "The device token to search for.",
				},
				"active": &graphql.ArgumentConfig{
					Type:        graphql.Boolean,
					Description: "The state of a device associated with this user",
				},
				"intersection": &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					Description:  "If the search arguments should be accepted in union or intersection.",
					DefaultValue: false,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userLogin := p.Args["login"].(string)
				result := User{}
				query := bson.M{"$or": []bson.M{{"email": userLogin}, {"number": userLogin}}}
				response := bson.M{"services": 1}

				searchParams := make(map[string]interface{})

				if p.Args["name"] != nil {
					searchParams["name"] = p.Args["name"].(string)
				}

				if p.Args["token"] != nil {
					searchParams["token"] = p.Args["token"].(string)
				}

				if p.Args["active"] != nil {
					searchParams["active"] = p.Args["active"].(bool)
				}

				if len(searchParams) > 0 {
					if p.Args["intersection"] == false {
						params := make([]bson.M, 1)
						for e, r := range searchParams {
							params = append(params, bson.M{e: r})
						}
						query["services"] = bson.M{"$elemMatch": bson.M{"$or": params}}
					} else {
						query["services"] = bson.M{"$elemMatch": searchParams}
					}
				}

				err := users.Find(query).Select(response).One(&result)
				if err != nil {
					return nil, err
				}
				return result.Devices, nil
			},
		},
	}
	mutations := graphql.Fields{
		"addUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"user": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(userInputType),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				val := p.Args["user"].(map[string]interface{})

				if val["number"] == nil && val["email"] == nil {
					return nil, errors.New("Either an email or phone number is required")
				}
				user := &User{}

				if val["email"] != nil {
					user.Email = val["email"].(string)
				}
				if val["number"] != nil {
					user.Number = val["number"].(string)
				}

				query := bson.M{"email": user.Email}
				co, err := users.Find(query).Count()
				if err != nil {
					return nil, err
				}
				if co > 0 {
					return nil, errors.New("A user already exists with that email")
				}

				query = bson.M{"number": user.Number}
				co, err = users.Find(query).Count()
				if err != nil {
					return nil, err
				}
				if co > 0 {
					return nil, errors.New("A user already exists with that number")
				}

				err = users.Insert(user)
				if err != nil {
					return nil, err
				}
				return user, nil
			},
		},
		"addService": &graphql.Field{
			Type: serviceType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"type": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"login": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name := p.Args["name"].(string)
				serviceProvider := p.Args["type"].(string)
				login := p.Args["login"].(string)

				service := &Service{Name: name, Type: serviceProvider, Active: true}
				query := bson.M{"$or": []bson.M{{"email": login}, {"number": login}}}
				update := bson.M{"$push": bson.M{"services": service}}
				err := users.Update(query, update)
				if err != nil {
					return nil, err
				}
				return service, nil
			},
		},

		"addDevice": &graphql.Field{
			Type: deviceType,
			Args: graphql.FieldConfigArgument{
				"device": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(deviceInput),
				},
				"login": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				device := p.Args["device"].(map[string]interface{})
				userLogin := p.Args["login"].(string)

				query := bson.M{"email": userLogin}
				update := bson.M{"$push": bson.M{"devices": &Device{Name: device["name"].(string), Token: device["token"].(string), Active: device["active"].(bool)}}}
				err := users.Update(query, update)
				if err != nil {
					return nil, err
				}
				return device, nil
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

	fmt.Println(config.Mongo.Hostname)
	session := newDatabaseConnection(config.Mongo.Hostname)
	session.SetMode(mgo.Monotonic, true)

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

func newDatabaseConnection(hostname string) (session *mgo.Session) {
	session, err := mgo.Dial(hostname)
	if err != nil {
		log.Fatal(err)
	}
	return
}
