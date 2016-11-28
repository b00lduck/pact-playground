package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/Sirupsen/logrus"
	"github.com/pact-foundation/pact-go/types"
)

var dir, _ = os.Getwd()
var pactDir = fmt.Sprintf("%s/../../pacts", dir)
var logDir = fmt.Sprintf("%s/log", dir)
var pact1 dsl.Pact
var pact2 dsl.Pact
var req *http.Request

const pactDaemonPort = 6666

func TestMain(m *testing.M) {

	setup()

	code := m.Run()

	// Write
	pact1.WritePact()

	// Publish the PACT to the PACT broker
	pr := types.PublishRequest{
		PactBroker:             "http://localhost",
		PactURLs:               []string{"../../pacts/myconsumer-myprovider.json"},
		ConsumerVersion:        "1.0.0",
		Tags:                   []string{"latest", "dev"},
	}
	pb := dsl.Publisher{}
	err := pb.Publish(pr)
	if err != nil {
		logrus.Error(err)
	}

	pact1.Teardown()
	pact2.Teardown()

	os.Exit(code)
}

// Setup common test data
func setup() {
	pact1 = createPact1()
	pact2 = createPact2()

	// Create a request to pass to our handler.
	req, _ = http.NewRequest("POST", "/login", strings.NewReader(""))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
}

func createPact1() dsl.Pact {
	pact := dsl.Pact{
		Port:     pactDaemonPort,
		Consumer: "MyConsumer",
		Provider: "MyProvider",
		LogDir:   logDir,
		PactDir:  pactDir,
	}
	pact.Setup()
	return pact
}

func createPact2() dsl.Pact {
	pact := dsl.Pact{
		Port:     pactDaemonPort,
		Consumer: "MyConsumer",
		Provider: "MyProvider",
		LogDir:   logDir,
		PactDir:  pactDir,
	}
	pact.Setup()
	return pact
}

func TestPactConsumerLoginHandler_UserUnauthorised(t *testing.T) {

	var testAccessProvider = func() error {
		err := AccessProvider(fmt.Sprintf("http://localhost:%d", pact1.Server.Port))
		if err != nil {
			return fmt.Errorf("Error occured: %v", err)
		}
		return nil
	}

	pact1.
	AddInteraction().
		UponReceiving("Some test request").
		WithRequest(dsl.Request{
		Method: "GET",
		Path:   "/test.json",
	}).
		WillRespondWith(dsl.Response{
		Status: 200,
		Headers: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
		Body: `{"foo":"bar"}`,
	})

	err := pact1.Verify(testAccessProvider)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}

}
