# Lovi

[![Build Status](https://travis-ci.org/pakerfeldt/lovi.svg?branch=master)](https://travis-ci.org/pakerfeldt/lovi)

Lovi is a lightweight pager application for message distribution. It offers a web api for triggering events and uses a policy configuration to determine how and to whom the events gets distributed.

This application currently supports the following transport mechanisms:

- SMS / Phone calls (through 46elk.se)
- stdout (for logging)

## How does it work?

Lovi is configured using a yaml file. Here you activate the transports you want to use and define the set of policies for which you plan to send events. A policy describes how events will be sent and to whom. See [Configuration](https://github.com/pakerfeldt/lovi/wiki/Configuration) wiki page for examples.

You can have multiple policies for different needs, each with its own set of configuration.

## Contribute

See [Transports](https://github.com/pakerfeldt/lovi/wiki/Transports) wiki page for a guide on how to add new transports.
