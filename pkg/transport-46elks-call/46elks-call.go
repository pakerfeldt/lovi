package transport46elksCall

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pakerfeldt/lovi/pkg/transports"
)

type Call struct {
	ID string
}

type Call46Elks struct {
	ack             transports.Acknowledge
	username        string
	password        string
	sender          string
	ackSuffix       string
	baseURLIncoming string
	text2speech     string
	calls           map[string]string
}

func Id() string {
	return "46elks-call"
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
	ackSuffix, exists := config["ackSuffix"]
	baseURLIncoming, exists := config["baseUrlIncoming"]
	if !exists {
		baseURLIncoming = ""
	}
	text2speech, exists := config["textToSpeech"]
	if !exists {
		text2speech = ""
	}
	transport := Call46Elks{ack: ack,
		username:        username,
		password:        password,
		sender:          sender,
		ackSuffix:       ackSuffix,
		baseURLIncoming: baseURLIncoming,
		text2speech:     text2speech,
		calls:           make(map[string]string)}
	router.HandleFunc("/transports/46elks/callaction", transport.incomingCallAction).Methods("POST")
	return transport
}

func (Call46Elks) Id() string {
	return Id()
}

func (c Call46Elks) incomingCallAction(w http.ResponseWriter, r *http.Request) {
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

func (c Call46Elks) Send(id string, message string, target string, ack bool) {
	if ack {
		message += " " + c.ackSuffix
	}
	sound := "sound/beep"
	if c.text2speech != "" {
		encodedText := base64.StdEncoding.EncodeToString([]byte(message))
		sound = strings.Replace(c.text2speech, "{base64}", encodedText, -1)
	}

	next := ""
	if c.baseURLIncoming != "" {
		next = ", \"next\": \"" + c.baseURLIncoming + "/transports/46elks/callaction" + "\""
	}

	data := url.Values{
		"from":        {c.sender},
		"to":          {target},
		"voice_start": {"{\"ivr\": \"" + sound + "\",\"digits\": 1,\"timeout\": 10,\"repeat\": 3" + next + "}"}}

	log.Printf("Calling %s with message %s\n", target, message)

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

	if resp.StatusCode != 200 {
		log.Printf("46elk-call: Invalid response code %d, '%s'\n", resp.StatusCode, body)
	} else {
		var call Call
		err = json.Unmarshal([]byte(body), &call)
		if err != nil {
			log.Printf("46elk-call: Could not parse call response %s, err was %v\n", body, err)
		}
		c.calls[call.ID] = id
	}

}
