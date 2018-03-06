package main

// Service is the database structure of service
type Service struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Active      bool   `json:"active"`
	Description string `json:"description"`
	URI         string `json:"uri"`
}

// Device is the database structure of device
type Device struct {
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	Token       string `json:"token"`
	Description string `json:"description"`
}

// User datatype
type User struct {
	Email    string    `json:"email"`
	Number   string    `json:"number"`
	Status   string    `json:"status"`
	Services []Service `json:"services"`
	Devices  []Device  `json:"devices"`
}

// Login datatype
type Login struct {
	User        User   `json:"user"`
	LoginMethod string `json:"loginMethod"`
}

type TokenValidation struct {
	Email           string `json:"email"`
	Number          string `json:"number"`
	SMSCode         string `json:"SMSCode"`
	EmailToken      string `json:"emailToken"`
	TokenHardExpire int64  `json:"tokenHardExpire"`
}
