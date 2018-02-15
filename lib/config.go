package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config structure, shows all options.
type Config struct {
	Port  string `json:"port"`
	Mongo struct {
		Hostname string `json:"hostname"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"mongo"`
	SMTP   SMTPConfig   `json:"smtp"`
	Twilio TwilioConfig `json:"twilio"`
}

// SMTPConfig contains the server configuration to send email notifications
type SMTPConfig struct {
	Server        string `json:"server"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Port          string `json:"port"`
	SenderAddress string `json:"senderAddress"`
}

// TwilioConfig contains the proper authorative credentials to contact the Twilio API.
type TwilioConfig struct {
	AccountSid   string `json:"accountSID"`
	AuthToken    string `json:"authToken"`
	SenderNumber string `json:"senderNumber"`
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
