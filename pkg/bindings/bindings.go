package bindings

import (
	"github.com/gorilla/mux"
	"github.com/pakerfeldt/lovi/pkg/transport-46elks-call"
	"github.com/pakerfeldt/lovi/pkg/transport-46elks-sms"
	"github.com/pakerfeldt/lovi/pkg/transport-log"
	"github.com/pakerfeldt/lovi/pkg/transports"
)

var factories map[string]func(*mux.Router, map[string]string, transports.Acknowledge) transports.Transport

func init() {
	factories = make(map[string]func(*mux.Router, map[string]string, transports.Acknowledge) transports.Transport)
	factories[transportLog.Id()] = transportLog.CreateTransport
	factories[transport46elksSms.Id()] = transport46elksSms.CreateTransport
	factories[transport46elksCall.Id()] = transport46elksCall.CreateTransport
}

func Factories() map[string](func(*mux.Router, map[string]string, transports.Acknowledge) transports.Transport) {
	return factories
}
