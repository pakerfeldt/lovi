package models

import "time"

type Config struct {
	Policies   []Policy
	Transports []TransportConfig
}

type Settings struct {
	ConfigFile string
	Port       int
}

type TransportConfig struct {
	ID            string
	Configuration map[string]string
}

type Event struct {
	ID            string
	Time          time.Time
	Message       string
	Policy        string
	ReceiverIndex int
	AckedAt       time.Time
	AckedBy       string
}

func (e Event) Id() string {
	return e.ID
}

func (e Event) Acked(ackedAt time.Time, ackedBy string) Event {
	return Event{ID: e.ID,
		Time:          e.Time,
		Message:       e.Message,
		Policy:        e.Policy,
		ReceiverIndex: e.ReceiverIndex,
		AckedAt:       ackedAt,
		AckedBy:       ackedBy}
}

type Policy struct {
	ID                      string
	Ack                     bool
	AckTimeoutSeconds       int `yaml:"ackTimeoutSeconds,omitempty"`
	AutoResolveAfterSeconds int `yaml:"autoResolveAfterSeconds,omitempty"`
	Broadcast               bool
	Receivers               []Receiver
}

func (p Policy) Id() string {
	return p.ID
}

type Receiver struct {
	Name      string
	Transport string
	Target    string
}
