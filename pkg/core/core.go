package core

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/pakerfeldt/lovi/pkg/idgenerator"
	"github.com/pakerfeldt/lovi/pkg/models"
	t "github.com/pakerfeldt/lovi/pkg/transports"
)

var transports map[string]t.Transport
var policies map[string]models.Policy
var unAckedEvents *UnAckedEventsMap

func AddPolicy(policy models.Policy) {
	_, exists := policies[policy.ID]
	if exists {
		panic(errors.New("Policy '" + policy.ID + "' already exists."))
	}
	policies[policy.ID] = policy
}

func AddTransport(transport t.Transport) {
	_, exists := transports[transport.Id()]
	if exists {
		panic(errors.New("Transport '" + transport.Id() + "' already exists."))
	}
	transports[transport.Id()] = transport
}

func initialise() {
	policies = make(map[string]models.Policy)
	transports = make(map[string]t.Transport)
	unAckedEvents = NewUnAckedEventsMap()
}

func init() {
	initialise()
}

func HandleAck(ID string) {
	log.Printf("Event '%s' acked", ID)
	unAckedEvents.Delete(ID)
}

func HandleEvent(policyId string, message string) error {
	id := idgenerator.New()
	messageReplaced := strings.Replace(message, "{{ID}}", id, -1)
	event := models.Event{ID: id, Time: time.Now(), Policy: policyId, Message: messageReplaced}
	policy, ok := policies[event.Policy]
	if !ok {
		return errors.New("Unknown policy: " + event.Policy)
	}
	if policy.Ack {
		unAckedEvents.Store(event.ID, event)
		scheduleEventWithAckRequired(time.NewTimer(0), event, policy)
	} else {
		sendEvent(event, policy, false)
	}
	return nil
}

func sendEvent(event models.Event, policy models.Policy, ackRequired bool) {
	if policy.Broadcast {
		for _, receiver := range policy.Receivers {
			transport := transports[receiver.Transport]
			target := receiver.Target
			go transport.Send(event.ID, event.Message, target, ackRequired)
		}
	} else {
		transport := transports[policy.Receivers[event.ReceiverIndex].Transport]
		target := policy.Receivers[event.ReceiverIndex].Target
		go transport.Send(event.ID, event.Message, target, ackRequired)
	}
}

func scheduleEventWithAckRequired(timer *time.Timer, event models.Event, policy models.Policy) {
	go func() {
		<-timer.C
		_, notAcked := unAckedEvents.Load(event.ID)
		if notAcked {
			log.Printf("Event not acked - triggered: %s\n", event.Time.String())
			if policy.AutoResolveAfterSeconds > 0 && time.Now().Sub(event.Time).Seconds() > float64(policy.AutoResolveAfterSeconds) {
				log.Printf("Event auto resolved after %d seconds\n", policy.AutoResolveAfterSeconds)
				unAckedEvents.Delete(event.ID)
			} else {
				sendEvent(event, policy, true)
				scheduleEventWithAckRequired(time.NewTimer(time.Duration(policy.AckTimeoutSeconds)*time.Second), event, policy)
			}
		}
	}()
}
