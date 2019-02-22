package models

import "time"

type Config struct {
	Policies   []Policy
	Transports []TransportConfig
}

type TransportConfig struct {
	ID            string
	Configuration map[string]string
}

// Alert represents an outgoing alert
type Alert struct {
	ID            string
	Time          time.Time
	Message       string
	Policy        string
	ReceiverIndex int
	AckedAt       time.Time
	AckedBy       string
}

func (a Alert) Id() string {
	return a.ID
}

func (a Alert) Acked(ackedAt time.Time, ackedBy string) Alert {
	return Alert{ID: a.ID,
		Time:          a.Time,
		Message:       a.Message,
		Policy:        a.Policy,
		ReceiverIndex: a.ReceiverIndex,
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
