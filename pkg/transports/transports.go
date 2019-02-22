package transports

type Acknowledge func(string)

type Transport interface {
	Id() string
	Send(id string, message string, target string, ack bool)
}
