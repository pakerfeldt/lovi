package transportLog

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/pakerfeldt/lovi/pkg/transports"
)

func CreateTransport(router *mux.Router, config map[string]string, ack transports.Acknowledge) transports.Transport {
	return TransportLog{}
}

func Id() string {
	return "log"
}

type TransportLog struct{}

func (TransportLog) Id() string {
	return Id()
}

func (TransportLog) Send(id string, message string, target string, ackRequired bool) {
	log.Print("Message (" + id + ") to " + target + ": " + message)
}
