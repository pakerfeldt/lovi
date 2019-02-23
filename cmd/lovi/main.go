package main

import (
	"errors"

	"github.com/gorilla/mux"
	"github.com/pakerfeldt/lovi/pkg/api"
	"github.com/pakerfeldt/lovi/pkg/bindings"
	"github.com/pakerfeldt/lovi/pkg/config"
	"github.com/pakerfeldt/lovi/pkg/core"
)

func main() {
	settings := config.Settings()
	router := mux.NewRouter()

	conf := config.Parse(settings.ConfigFile)
	for _, policy := range conf.Policies {
		core.AddPolicy(policy)
	}

	for _, transport := range conf.Transports {
		factory, exists := bindings.Factories()[transport.ID]
		if !exists {
			panic(errors.New("Transport '" + transport.ID + "' does not exist."))
		}
		core.AddTransport(factory(router, transport.Configuration, core.HandleAck))
	}

	api.Init(router, settings)
}
