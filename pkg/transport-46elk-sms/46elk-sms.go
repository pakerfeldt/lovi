package transport46elkSms

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pakerfeldt/lovi/pkg/transports"
)

type Sms46Elk struct {
	ack      transports.Acknowledge
	username string
	password string
	sender   string
}

func Id() string {
	return "46elk-sms"
}

func CreateTransport(router *mux.Router, config map[string]string, ack transports.Acknowledge) transports.Transport {
	sender, exists := config["sender"]
	if !exists {
		sender = "Lovi"
	}
	username, exists := config["username"]
	if !exists {
		panic(errors.New(Id() + " requires 'username' in configuration."))
	}
	password, exists := config["password"]
	if !exists {
		panic(errors.New(Id() + " requires 'password' in configuration."))
	}
	transport := Sms46Elk{ack: ack, username: username, password: password, sender: sender}
	router.HandleFunc("/transports/46elks/sms", transport.incomingSms).Methods("POST")
	return transport
}

func (Sms46Elk) Id() string {
	return Id()
}

func (s Sms46Elk) incomingSms(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")
	s.ack(message)
}

func (s Sms46Elk) Send(id string, message string, target string, ack bool) {
	data := url.Values{
		"from":    {s.sender},
		"to":      {target},
		"message": {message}}

	req, err := http.NewRequest("POST", "https://api.46elks.com/a1/SMS", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.SetBasicAuth(s.username, s.password)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("46elk-sms: Error sending sms: %s", err)
	}

	defer resp.Body.Close()
}
