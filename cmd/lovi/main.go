package main

import (
	"errors"
	"flag"

	"github.com/gorilla/mux"
	"github.com/pakerfeldt/lovi/pkg/api"
	"github.com/pakerfeldt/lovi/pkg/bindings"
	"github.com/pakerfeldt/lovi/pkg/config"
	"github.com/pakerfeldt/lovi/pkg/core"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.Parse()

	if configFile == "" {
		panic(errors.New("No configuration file specified"))
	}
	router := mux.NewRouter()
	conf := config.Parse(configFile)
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
	config.Print(conf)
	api.Init(router)
}
