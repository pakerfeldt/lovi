# Lovi

Lovi is a light weight pager application for message distribution. Use the web API to trigger events that gets sent to receivers according to a pre-defined policy.

This application currently supports the following transport mechanisms:

- SMS / Phone calls (through 46elk.se)
- stdout (for logging)

## How does it work?

Lovi is configured using a yaml files. Here you define the set of policies that suits your need. A policy defines how a message gets distributed and contains an ordered list of receivers to receive the events. You can choose to require an acknowledgement and have the message be sent to each receiver until someone acknowledge the message.

You can have multiple policies for different needs, each with its own set of configuration.

## Extending

Writing a new transport is fairly easy. Add a new `transport-xyz` under `pkg`. Then make sure to provide a constructor method on the form `func(*mux.Router, map[string]string, transports.Acknowledge)transports.Transport`. The function takes a `mux.Router` in case the transport needs to listen to incoming web requests. Just add a handler to the router, nothing else. E.g. `router.HandleFunc("/transports/xyz/my/custom/path", myHandler).Methods("GET")`. The map contains any configuration provided through the yaml file. Lastly, the `transports.Acknowledge` is a function to call when the transport receives acknowledgements.
