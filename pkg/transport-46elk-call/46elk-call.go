package transport46elkCall

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pakerfeldt/lovi/pkg/transports"
)

type Call struct {
	ID string
}

type Call46Elk struct {
	ack      transports.Acknowledge
	username string
	password string
	sender   string
	calls    map[string]string
}

func Id() string {
	return "46elk-call"
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
	transport := Call46Elk{ack: ack,
		username: username,
		password: password,
		sender:   sender,
		calls:    make(map[string]string)}
	router.HandleFunc("/transports/46elks/callaction", transport.incomingCallAction).Methods("POST")
	return transport
}

func (Call46Elk) Id() string {
	return Id()
}

func (c Call46Elk) incomingCallAction(w http.ResponseWriter, r *http.Request) {
	result := r.FormValue("result")
	callId := r.FormValue("callid")
	if result == "1" {
		eventId, exists := c.calls[callId]
		if exists {
			delete(c.calls, callId)
			c.ack(eventId)
		} else {
			log.Printf("46elk-call: Unknown event associated with call id %s\n", callId)
		}
	}
	fmt.Fprintf(w, "{\"hangup\": \"busy\"}")
}

func (c Call46Elk) Send(id string, message string, target string, ack bool) {
	data := url.Values{
		"from":        {c.sender},
		"to":          {target},
		"voice_start": {"{\"ivr\": \"sound/beep\",\"digits\": 1,\"timeout\": 10,\"repeat\": 3,\"next\": \"http://node.akerfeldt.se:8080/transports/46elks/callaction\"}"}}

	req, err := http.NewRequest("POST", "https://api.46elks.com/a1/Calls", bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.SetBasicAuth(c.username, c.password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("46elk-call: Error making call: %v\n", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("46elk-call: Error reading response: %v\n", err)
	}

	var call Call
	err = json.Unmarshal([]byte(body), &call)
	if err != nil {
		log.Printf("46elk-call: Could not parse call response, err was %v\n", err)
	}
	c.calls[call.ID] = id
}
