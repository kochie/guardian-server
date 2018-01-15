package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//Config structure, shows all options.
type Config struct {
	Port          string `json:"port"`
	MongoHostname string `json:"mongo_hostname"`
}

//ImportConfig will import the configuration options from the config file.
func ImportConfig() (conf *Config) {
	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(raw, &conf)

	return
}
