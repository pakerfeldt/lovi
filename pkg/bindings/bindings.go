package bindings

import (
	"github.com/gorilla/mux"
	"github.com/pakerfeldt/lovi/pkg/transport-46elk-call"
	"github.com/pakerfeldt/lovi/pkg/transport-46elk-sms"
	"github.com/pakerfeldt/lovi/pkg/transport-log"
	"github.com/pakerfeldt/lovi/pkg/transports"
)

var factories map[string]func(*mux.Router, map[string]string, transports.Acknowledge) transports.Transport

func init() {
	factories = make(map[string]func(*mux.Router, map[string]string, transports.Acknowledge) transports.Transport)
	factories[transportLog.Id()] = transportLog.CreateTransport
	factories[transport46elkSms.Id()] = transport46elkSms.CreateTransport
	factories[transport46elkCall.Id()] = transport46elkCall.CreateTransport
}

func Factories() map[string](func(*mux.Router, map[string]string, transports.Acknowledge) transports.Transport) {
	return factories
}
