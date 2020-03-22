// +build acceptance

package eventing

import (
	"fmt"
	"testing"

	tester "github.com/kyma-project/kyma/tests/console-backend-service"
	"github.com/kyma-project/kyma/tests/console-backend-service/internal/graphql"
	"github.com/stretchr/testify/assert"
)

const (
	TriggerName          = "TestTrigger"
	TriggerNamespace     = "kyma-system"
	SubscriberName       = "TestService"
	SubscriberNamespace  = "kyma-system"
	SubscriberAPIVersion = "eventing.knative.dev/v1alpha1"
	SubscriberKind       = "Trigger"
	BrokerName           = "default"
)

func TestTriggerEventQueries(t *testing.T) {
	c, err := graphql.New()
	assert.NoError(t, err)

	//eventingCli, _, err := client.NewDynamicClientWithConfig()
	//require.NoError(t, err)

	//Subscribe events
	subscription := subscribeTriggerEvent(c, triggerArgumentFields(""), triggerEventDetailsFields())
	defer subscription.Close()

	//Create Trigger
	err = mutationTrigger(c,"create", createTriggerArguments(), triggerDetailsFields())
	assert.NoError(t, err)

	//Check and compare events
	event, err := readTriggerEvent(subscription)
	assert.NoError(t, err)

	expectedEvent := newTriggerEvent("ADD", fixTrigger())
	checkTriggerEvent(t, expectedEvent, event)

	//List triggers
	err = listTriggers(c, listTriggersArguments(), triggerDetailsFields())
	assert.NoError(t, err)
}

func newTriggerEvent(eventType string, trigger Trigger) TriggerEvent {
	return TriggerEvent{
		Type:    eventType,
		Trigger: trigger,
	}
}

func checkTriggerEvent(t *testing.T, expected, actual TriggerEvent) {
	assert.Equal(t, expected.Type, actual.Type)
	assert.Equal(t, expected.Trigger.Name, actual.Trigger.Name)
	assert.Equal(t, expected.Trigger.Namespace, actual.Trigger.Namespace)
}

func checkTriggerList(t *testing.T, triggers []Trigger) {
	assert.Equal(t, )
	assert.Equal(t, trigger.Namespace, TriggerNamespace)
	assert.Equal(t, trigger.Namespace, TriggerNamespace)
	assert.Equal(t, trigger.Namespace, TriggerNamespace)
}

func readTriggerEvent(sub *graphql.Subscription) (TriggerEvent, error) {
	type Response struct {
		TriggerEvent TriggerEvent
	}

	var response Response
	err := sub.Next(&response, tester.DefaultDeletionTimeout)

	return response.TriggerEvent, err
}

func listTriggers(client *graphql.Client, arguments, resourceDetailsQuery string) error {
	query := fmt.Sprintf(`
		query{
			triggers (
				%s
			){
				%s
			}
		}
	`)
	req := graphql.NewRequest(query)
	err := client.Do(req, nil)

	return err
}

func listTriggersArguments() string {
	return fmt.Sprintf(`
		namespace: "%s"
	`, TriggerNamespace)
}

func mutationTrigger(client *graphql.Client, requestType, arguments, resourceDetailsQuery string) error {
	query := fmt.Sprintf(`
		mutation {
			%sTrigger (
				%s
			){
				%s
			}
		}
	`,requestType, arguments, resourceDetailsQuery)
	req := graphql.NewRequest(query)
	err := client.Do(req, nil)

	return err
}

func createTriggerArguments() string {
	return fmt.Sprintf(`
		trigger: {
			name: "%s",
			namespace: "%s",
			broker: "%s"
			subscriber: {
				ref: {
					apiVersion: "%s",
					kind: "%s",
					name: "%s",
					namespace: "%s"
				}
			}
		},
	`, TriggerName, TriggerNamespace, BrokerName, SubscriberAPIVersion, SubscriberKind, SubscriberName, SubscriberNamespace)
}

func subscribeTriggerEvent(client *graphql.Client, arguments, resourceDetailsQuery string) *graphql.Subscription {
	query := fmt.Sprintf(`
		subscription {
			triggerEvent (
				%s
			){
				%s
			}
		}
	`, arguments, resourceDetailsQuery)
	req := graphql.NewRequest(query)

	return client.Subscribe(req)
}

func triggerArgumentFields(namespace string) string {
	return fmt.Sprintf(`
		namespace: "%s",
		subscriber: {
			ref: {
				apiVersion: "%s",
				kind: "%s",
				name: "%s",
				namespace: "%s"
			}
		}
	`, namespace, SubscriberAPIVersion, SubscriberKind, SubscriberName, SubscriberNamespace)
}

func triggerDetailsFields() string {
	return `
		name
    	namespace
		broker
    	filterAttributes
		subscriber {
			uri
			ref {
				apiVersion
				kind
				name
				namespace
			}
		}
    	status {
		reason
		status
		}
	`
}

func triggerEventDetailsFields() string {
	return fmt.Sprintf(`
        type
        trigger {
			%s
        }
    `, triggerDetailsFields())
}

func fixTrigger() Trigger {
	return Trigger{
		Name:      TriggerName,
		Namespace: TriggerNamespace,
	}
}