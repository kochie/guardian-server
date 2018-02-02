package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func checkForUser(users *mgo.Collection, val map[string]interface{}) (interface{}, error) {
	user := User{}

	if val["number"] == nil && val["email"] == nil {
		return nil, errors.New("Either an email or phone number is required")
	}

	if val["email"] != nil {
		user.Email = val["email"].(string)
		query := bson.M{"email": user.Email}
		co, err := users.Find(query).Count()
		if err != nil {
			return nil, err
		}
		if co > 0 {
			return nil, errors.New("A user already exists with that email")
		}
	}

	if val["number"] != nil {
		user.Number = val["number"].(string)
		query := bson.M{"number": user.Number}
		co, err := users.Find(query).Count()
		if err != nil {
			return nil, err
		}
		if co > 0 {
			return nil, errors.New("A user already exists with that number")
		}
	}

	return user, nil
}

func addUser(users *mgo.Collection, val map[string]interface{}) (interface{}, error) {
	user, err := checkForUser(users, val)
	if err != nil {
		return nil, err
	}

	err = users.Insert(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func getServices(userLogin string, p graphql.ResolveParams, users *mgo.Collection) ([]Service, error) {
	result := User{}
	query := bson.M{"$or": []bson.M{{"email": userLogin}, {"number": userLogin}}}
	response := bson.M{"services": 1}

	if p.Args["filter"] != nil {
		filter := p.Args["filter"].(map[string]interface{})

		if len(filter) > 0 {
			if p.Args["intersection"] == false {
				params := make([]bson.M, 0)
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
		if err.Error() == "not found" {
			return []Service{}, nil
		}
		return nil, err

	}
	return result.Services, nil
}

func getDevices(userLogin string, p graphql.ResolveParams, users *mgo.Collection) ([]Device, error) {
	result := User{}
	query := bson.M{"$or": []bson.M{{"email": userLogin}, {"number": userLogin}}}
	response := bson.M{"devices": 1}

	if p.Args["filter"] != nil {
		filter := p.Args["filter"].(map[string]interface{})

		if len(filter) > 0 {
			if p.Args["intersection"] == false {
				params := make([]bson.M, 0)
				for e, r := range filter {
					params = append(params, bson.M{e: r})
				}
				query["devices"] = bson.M{"$elemMatch": bson.M{"$or": params}}
			} else {
				query["devices"] = bson.M{"$elemMatch": filter}
			}
		}
	}

	err := users.Find(query).Select(response).One(&result)
	if err != nil {
		if err.Error() == "not found" {
			return []Device{}, nil
		}
		return nil, err

	}
	return result.Devices, nil
}

// BuildQL will create a schema for graphql
func BuildQL(session *mgo.Session) {
	users := session.DB("guardian").C("users")

	serviceType := graphql.NewObject(graphql.ObjectConfig{
		Name: "ServiceType",
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
			"uri": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The location of the service",
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "An optional field to describe the device.",
			},
		},
	})

	deviceType := graphql.NewObject(graphql.ObjectConfig{
		Name: "DeviceType",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The name given to this device.",
			},
			"active": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.Boolean),
				Description: "The status of the device, if it is active or not.",
			},
			"token": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The JWT that authoriizes this device.",
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "An optional field to describe the device.",
			},
		},
	})

	serviceInput := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "Service",
			Fields: graphql.InputObjectConfigFieldMap{
				"name": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The name of a service associated with this user",
				},
				"type": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The service provider of a service associated with this user",
				},
				"active": &graphql.InputObjectFieldConfig{
					Type:         graphql.NewNonNull(graphql.Boolean),
					Description:  "The state of a service associated with this user",
					DefaultValue: true,
				},
				"description": &graphql.InputObjectFieldConfig{
					Type:         graphql.String,
					Description:  "An optional field to describe the service.",
					DefaultValue: "",
				},
				"uri": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The location of the service.",
				},
			},
		},
	)

	deviceInput := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "Device",
			Fields: graphql.InputObjectConfigFieldMap{
				"name": &graphql.InputObjectFieldConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The name of a device associated with this user",
				},
				"token": &graphql.InputObjectFieldConfig{
					Type:        graphql.String,
					Description: "The device token associated with this user",
				},
				"active": &graphql.InputObjectFieldConfig{
					Type:         graphql.Boolean,
					Description:  "The state of a service associated with this user",
					DefaultValue: true,
				},
				"description": &graphql.InputObjectFieldConfig{
					Type:        graphql.String,
					Description: "An optional field to describe the service.",
				},
			},
		},
	)

	userInput := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "User",
			Fields: graphql.InputObjectConfigFieldMap{
				"email": &graphql.InputObjectFieldConfig{
					Type:        graphql.String,
					Description: "Email registered to the user that can be used to identify the user.",
				},
				"number": &graphql.InputObjectFieldConfig{
					Type:        graphql.String,
					Description: "Phone number that can be used to identify the user.",
				},
			},
		},
	)

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "UserType",
		Fields: graphql.Fields{
			"email": &graphql.Field{
				Type:        graphql.String,
				Description: "Email address of the user. Used for email login links.",
			},
			"username": &graphql.Field{
				Type:        graphql.String,
				Description: "The username of the user. Can be used to search for other users.",
			},
			"number": &graphql.Field{
				Type:        graphql.String,
				Description: "The phone number associated with the user. This is the number SMS login codes are sent to.",
			},
			"status": &graphql.Field{
				Type:        graphql.String,
				Description: "The status of the account. Can be active|inactive|deactivated.",
			},
			"devices": &graphql.Field{
				Type:        &graphql.List{OfType: deviceType},
				Description: "A list of devices associated with the user.",
				Args: graphql.FieldConfigArgument{
					"filter": &graphql.ArgumentConfig{
						Type: deviceInput,
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
						return getDevices(user.Number, p, users)
					}
					return getDevices(user.Email, p, users)
				},
			},
			"services": &graphql.Field{
				Type:        &graphql.List{OfType: serviceType},
				Description: "A list of services associated with the user.",
				Args: graphql.FieldConfigArgument{
					"filter": &graphql.ArgumentConfig{
						Type: serviceInput,
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

	loginStatus := graphql.NewObject(graphql.ObjectConfig{
		Name: "LoginStatus",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type:        userType,
				Description: "The user information",
			},
			"loginMethod": &graphql.Field{
				Type:        graphql.String,
				Description: "The method of login used by the client.",
			},
		},
	})

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
					Type:        graphql.NewNonNull(graphql.String),
					Description: "The login to search for in the database.",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				login := p.Args["login"].(string)

				result := User{}
				searchFields := []bson.M{{"email": login}, {"number": login}, {"username": login}}
				if bson.IsObjectIdHex(login) {
					searchFields = append(searchFields, bson.M{"_id": bson.ObjectIdHex(login)})
				}
				query := bson.M{"$or": searchFields}
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
					Type:        serviceInput,
					Description: "Specify a filter to search services by",
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
				"filter": &graphql.ArgumentConfig{
					Type:        deviceInput,
					Description: "Specify a filter to search device by",
				},
				"intersection": &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					Description:  "If the search arguments should be accepted in union or intersection.",
					DefaultValue: false,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userLogin := p.Args["login"].(string)
				return getDevices(userLogin, p, users)
			},
		},
	}
	mutations := graphql.Fields{
		"addUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"user": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(userInput),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				val := p.Args["user"].(map[string]interface{})
				return addUser(users, val)
			},
		},

		"updateUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"login": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"user": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(userInput),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				login := p.Args["login"].(string)
				val := p.Args["user"].(map[string]interface{})
				user := User{}

				change := mgo.Change{
					Update:    bson.M{"$set": val},
					ReturnNew: true,
				}

				query := bson.M{"$or": []bson.M{{"email": login}, {"number": login}}}
				_, err := users.Find(query).Apply(change, &user)
				if err != nil {
					return nil, err
				}

				return user, nil
			},
		},

		"removeUser": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"login": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				login := p.Args["login"].(string)
				query := bson.M{"$or": []bson.M{{"email": login}, {"number": login}}}
				err := users.Remove(query)
				if err != nil {
					return nil, err
				}
				return true, nil
			},
		},

		"addService": &graphql.Field{
			Type: serviceType,
			Args: graphql.FieldConfigArgument{
				"service": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(serviceInput),
				},
				"login": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				service := p.Args["service"].(map[string]interface{})
				login := p.Args["login"].(string)

				serviceObject := &Service{
					Name:   service["name"].(string),
					Type:   service["type"].(string),
					Active: service["active"].(bool),
				}
				query := bson.M{"$or": []bson.M{{"email": login}, {"number": login}}}
				update := bson.M{"$push": bson.M{"services": serviceObject}}
				err := users.Update(query, update)
				if err != nil {
					return nil, err
				}
				return serviceObject, nil
			},
		},

		"updateService": &graphql.Field{
			Type: serviceType,
			Args: graphql.FieldConfigArgument{
				"service": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(serviceInput),
				},
				"updatedService": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(serviceInput),
				},
				"login": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				service := p.Args["service"].(map[string]interface{})
				login := p.Args["login"].(string)
				val := p.Args["updatedService"].(map[string]interface{})

				userObject := User{}

				change := bson.M{"$set": bson.M{"services.$": val}}

				selector := bson.M{"services.$": true}

				query := bson.M{"$or": []bson.M{{"email": login}, {"number": login}}, "services": bson.M{"$elemMatch": service}}
				err := users.Update(query, change)
				if err != nil {
					return nil, err
				}

				query = bson.M{"$or": []bson.M{{"email": login}, {"number": login}}, "services": bson.M{"$elemMatch": val}}
				err = users.Find(query).Select(selector).One(&userObject)
				if err != nil {
					return nil, err
				}

				return userObject.Services[0], nil
			},
		},

		"removeService": &graphql.Field{
			Type: serviceType,
			Args: graphql.FieldConfigArgument{
				"service": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(serviceInput),
				},
				"login": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				login := p.Args["login"].(string)
				service := p.Args["service"].(map[string]interface{})
				query := bson.M{"$or": []bson.M{{"email": login}, {"number": login}}}
				update := bson.M{"$pull": bson.M{"services": service}}
				err := users.Update(query, update)
				if err != nil {
					return nil, err
				}
				return true, nil
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
				login := p.Args["login"].(string)

				deviceObject := &Device{
					Name:   device["name"].(string),
					Token:  device["token"].(string),
					Active: device["active"].(bool),
				}
				query := bson.M{"$or": []bson.M{{"email": login}, {"number": login}}}
				update := bson.M{"$push": bson.M{"devices": deviceObject}}
				err := users.Update(query, update)
				if err != nil {
					return nil, err
				}
				return deviceObject, nil
			},
		},

		"updateDevice": &graphql.Field{
			Type: deviceType,
			Args: graphql.FieldConfigArgument{
				"device": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(deviceInput),
				},
				"updatedDevice": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(deviceInput),
				},
				"login": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				device := p.Args["device"].(map[string]interface{})
				login := p.Args["login"].(string)
				val := p.Args["updatedDevice"].(map[string]interface{})

				userObject := User{}

				change := bson.M{"$set": bson.M{"devices.$": val}}

				selector := bson.M{"devices.$": true}

				query := bson.M{"$or": []bson.M{{"email": login}, {"number": login}}, "devices": bson.M{"$elemMatch": device}}
				err := users.Update(query, change)
				if err != nil {
					return nil, err
				}

				query = bson.M{"$or": []bson.M{{"email": login}, {"number": login}}, "devices": bson.M{"$elemMatch": val}}
				err = users.Find(query).Select(selector).One(&userObject)
				if err != nil {
					return nil, err
				}

				return userObject.Devices[0], nil
			},
		},

		"removeDevice": &graphql.Field{
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
				login := p.Args["login"].(string)
				device := p.Args["device"].(map[string]interface{})
				query := bson.M{"$or": []bson.M{{"email": login}, {"number": login}}}
				update := bson.M{"$pull": bson.M{"devices": device}}
				err := users.Update(query, update)
				if err != nil {
					return nil, err
				}
				return true, nil
			},
		},

		"loginOrRegister": &graphql.Field{
			Type: loginStatus,
			Args: graphql.FieldConfigArgument{
				"user": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(userInput),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				val := p.Args["user"].(map[string]interface{})

				user, err := addUser(users, val)
				if err != nil {
					return nil, err
				}

				result := Login{
					User: user.(User),
				}

				if val["email"] != nil {
					result.LoginMethod = "email"
				}
				if val["number"] != nil {
					result.LoginMethod = "number"
				}

				return result, nil
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
