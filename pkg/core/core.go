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
var unAckedAlerts *UnAckedAlertsMap

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
	unAckedAlerts = NewUnAckedAlertsMap()
}

func init() {
	initialise()
}

func HandleAck(ID string) {
	log.Printf("Alert '%s' acked", ID)
	unAckedAlerts.Delete(ID)
}

func HandleAlert(policyId string, message string) error {
	id := idgenerator.New()
	messageReplaced := strings.Replace(message, "{{ID}}", id, -1)
	alert := models.Alert{ID: id, Time: time.Now(), Policy: policyId, Message: messageReplaced}
	policy, ok := policies[alert.Policy]
	if !ok {
		return errors.New("Unknown policy: " + alert.Policy)
	}
	if policy.Ack {
		unAckedAlerts.Store(alert.ID, alert)
		scheduleAlertWithAckRequired(time.NewTimer(0), alert, policy)
	} else {
		sendAlert(alert, policy, false)
	}
	return nil
}

func sendAlert(alert models.Alert, policy models.Policy, ackRequired bool) {
	if policy.Broadcast {
		for _, receiver := range policy.Receivers {
			transport := transports[receiver.Transport]
			target := receiver.Target
			go transport.Send(alert.ID, alert.Message, target, ackRequired)
		}
	} else {
		transport := transports[policy.Receivers[alert.ReceiverIndex].Transport]
		target := policy.Receivers[alert.ReceiverIndex].Target
		go transport.Send(alert.ID, alert.Message, target, ackRequired)
	}
}

func scheduleAlertWithAckRequired(timer *time.Timer, alert models.Alert, policy models.Policy) {
	go func() {
		<-timer.C
		_, notAcked := unAckedAlerts.Load(alert.ID)
		if notAcked {
			log.Printf("Alert not acked - triggered: %s\n", alert.Time.String())
			if policy.AutoResolveAfterSeconds > 0 && time.Now().Sub(alert.Time).Seconds() > float64(policy.AutoResolveAfterSeconds) {
				log.Printf("Alert auto resolved after %d seconds\n", policy.AutoResolveAfterSeconds)
				unAckedAlerts.Delete(alert.ID)
			} else {
				sendAlert(alert, policy, true)
				scheduleAlertWithAckRequired(time.NewTimer(time.Duration(policy.AckTimeoutSeconds)*time.Second), alert, policy)
			}
		}
	}()
}
