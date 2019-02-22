package core

import (
	"testing"

	"github.com/pakerfeldt/lovi/pkg/models"
	t "github.com/pakerfeldt/lovi/pkg/transports"
)

func TestAddPolicyTwiceCausesPanic(t *testing.T) {
	initialise()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Code did not panic when adding existing policy")
		}
	}()

	// The following is the code under test
	AddPolicy(models.Policy{ID: "someId"})
	AddPolicy(models.Policy{ID: "someId"})
}

func TestAddTransportTwiceCausesPanic(t *testing.T) {
	initialise()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Code did not panic when adding existing policy")
		}
	}()

	// The following is the code under test
	AddTransport(fakeTransport{})
	AddTransport(fakeTransport{})
}

func TestHandleAlertWithUnknownPolicyReturnsError(t *testing.T) {
	initialise()
	err := HandleAlert("unknownPolicy", "message")
	if err == nil {
		t.Error("Did not receive error handling alert with unknown policy")
	}
}

func TestHandleAlertWithoutAckDoesNotStoreAlertInUnAckedAlerts(t *testing.T) {
	initialise()
	AddPolicy(models.Policy{ID: "information", Ack: false,
		Receivers: []models.Receiver{models.Receiver{Transport: "fake", Target: "User1"}}})
	AddTransport(fakeTransport{})
	HandleAlert("information", "message")
	_, exists := unAckedAlerts.Load("someId")
	if exists {
		t.Error("Expecting non-ack alert to not be acknowledged")
	}
}

type fakeTransport struct{}

func (fakeTransport) Id() string {
	return "fake"
}

func (fakeTransport) Init(config map[string]string, ack t.Acknowledge) {

}

func (fakeTransport) Send(id string, message string, target string, ackRequired bool) {

}
