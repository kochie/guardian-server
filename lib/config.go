package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//Config structure, shows all options.
type Config struct {
	Port  string `json:"port"`
	Mongo struct {
		Hostname string `json:"hostname"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"mongo"`
	SMTP struct {
		Server   string `json:"server"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"smtp"`
}

//ImportConfig will import the configuration options from the config file.
func ImportConfig() (conf *Config) {
	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(raw, &conf)

	return conf
}
