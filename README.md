# Lovi

[![Build Status](https://travis-ci.org/pakerfeldt/lovi.svg?branch=master)](https://travis-ci.org/pakerfeldt/lovi)

Lovi is a lightweight pager application for message distribution. It offers a web api for triggering events and uses a policy configuration to determine how and to whom the events gets distributed.

This application currently supports the following transport mechanisms:

- SMS / Phone calls (through 46elk.se)
- stdout (for logging)

## How does it work?

Lovi is configured using a yaml file. Here you activate the transports you want to use and define the set of policies for which you plan to send events. A policy describes how events will be sent and to whom. See [Configuration](https://github.com/pakerfeldt/lovi/wiki/Configuration) wiki page for examples.

You can have multiple policies for different needs, each with its own set of configuration.

## Running

The recommended way of running lovi is through Docker.
`docker run -i -t pakerfeldt/lovi:1.0.0__linux_amd64`.
Lovi will try to read /config.yaml and listen to port 8080. You may want to change this by setting the `CONFIG` and `PORT` environment variables respectively.

## Contribute

See [Transports](https://github.com/pakerfeldt/lovi/wiki/Transports) wiki page for a guide on how to add new transports.
