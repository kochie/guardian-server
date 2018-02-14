package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"github.com/kochie/guardian-server/lib"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func validateUserLogin() {

}

// authoriseSMS should create and send an SMS code to the selected phone number and upload the confirmation code to the database so verifySMS can find it.
func authoriseSMS(phoneNumber string, config lib.TwilioConfig, validations *mgo.Collection) (bool, error) {
	token := GenerateToken(5)

	expireTime := time.Now().Add(time.Minute)
	err := validations.Insert(bson.M{"number": phoneNumber, "SMSCode": token, "tokenHardExpire": expireTime})
	if err != nil {
		return false, err
	}

	// Set initial variables
	accountSid := config.AccountSid
	authToken := config.AuthToken
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	// Build out the data for our message
	v := url.Values{}
	v.Set("To", phoneNumber)
	v.Set("From", "Guardian")
	v.Set("Body", token+" is your passcode.")
	v.Set("FriendlyName", "Guardian")
	rb := *strings.NewReader(v.Encode())

	// Create client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	fmt.Println(resp.Status)
	return true, nil
}

func authoriseEmail(emailAddress string, config lib.SMTPConfig, validations *mgo.Collection) (bool, error) {
	token := GenerateToken(15)
	err := validations.Insert(bson.M{"email": emailAddress, "emailToken": token})
	if err != nil {
		return false, err
	}
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Server)
	to := []string{emailAddress}
	msg := []byte("To: " + emailAddress + "\r\n" +
		"Subject: Guardian Magic Link\r\n" +
		"From: " + config.SenderAddress + "\r\n" +
		"\r\n" +
		"Hello, \r\n\r\nthis is your magic link to log into guardian: " + token + "\r\n")
	err = smtp.SendMail(config.Server+":"+config.Port, auth, config.SenderAddress, to, msg)
	if err != nil {
		return false, err
	}

	return true, nil
}

func validateEmail(user User, emailToken string, validations *mgo.Collection) (bool, error) {
	// Look for the token in the database, if found
	validation := TokenValidation{}
	query := bson.M{"emailToken": emailToken}
	if user.Email != "" {
		query["email"] = user.Email
	}
	if user.Number != "" {
		query["number"] = user.Number
	}

	change := mgo.Change{Remove: true}

	info, err := validations.Find(query).Apply(change, &validation)
	if err != nil {
		return false, err
	}

	if info.Removed == 1 {
		return true, nil
	}
	return false, nil

}

func validateSMS(user User, smsCode string, validations *mgo.Collection) (bool, error) {
	validation := TokenValidation{}
	query := bson.M{"SMSCode": smsCode}
	if user.Email != "" {
		query["email"] = user.Email
	}
	if user.Number != "" {
		query["number"] = user.Number
	}

	change := mgo.Change{Remove: true}

	info, err := validations.Find(query).Apply(change, &validation)
	if err != nil {
		return false, err
	}

	if info.Removed == 1 {
		return true, nil
	}
	return false, nil
}

func verifyUserLogin(user User, config lib.Config, validations *mgo.Collection) (bool, error) {
	if user.Email != "" {
		return authoriseEmail(user.Email, config.SMTP, validations)
	}
	if user.Number != "" {
		return authoriseSMS(user.Number, config.Twilio, validations)
	}
	return false, errors.New("No valid phone number or email address was given")

}

// Should verfiy login credentials in User object and return a device token
func login(user User, device Device, config lib.Config, validations *mgo.Collection) error {
	// should make a call to validate
	if ok, err := verifyUserLogin(user, config, validations); !ok {
		return err
	}
	return nil

}

// should log user out, invalidate their device token.
func logout(user User, deviceToken string, users *mgo.Collection) (bool, error) {
	deviceUser := User{}
	change := mgo.Change{Remove: true}
	query := bson.M{"devices": bson.M{"$elemMatch": bson.M{"token": deviceToken}}}
	if user.Email != "" {
		query["email"] = user.Email
	}
	if user.Number != "" {
		query["number"] = user.Number
	}
	info, err := users.Find(query).Apply(change, &deviceUser)

	if err != nil {
		return false, err
	}

	if info.Removed == 1 {
		return true, nil
	}
	return false, nil
}

// should create user account, log them in, and return a device token.
func register(user User, users *mgo.Collection, validations *mgo.Collection, config lib.Config, device Device) error {
	// add user to database in pre-registered state
	// validate the users login
	// either send a magic email, number code (in future can send a QR code and key for 2FA)
	err := addUser(users, user)
	if err != nil {
		return err
	}
	err = login(user, device, config, validations)
	if err != nil {
		return err
	}
	return nil
}
