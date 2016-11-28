package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"
	"github.com/Sirupsen/logrus"
)

var dir, _ = os.Getwd()
var pactDir = fmt.Sprintf("%s/../../pacts", dir)
var logDir = fmt.Sprintf("%s/log", dir)
var pact dsl.Pact
var rr http.ResponseWriter
var req *http.Request

// Use this to control the setup and teardown of Pact
func TestMain(m *testing.M) {
	// Setup Pact and related test stuff
	setup()

	// Run all the tests
	code := m.Run()

	// Shutdown the Mock Service and Write pact files to disk

	pact.WritePact()

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


	pact.Teardown()

	os.Exit(code)
}

// Setup common test data
func setup() {
	pact = createPact()

	// Create a request to pass to our handler.
	req, _ = http.NewRequest("POST", "/login", strings.NewReader(""))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Record response (satisfies http.ResponseWriter)
	rr = httptest.NewRecorder()
}

// Create Pact connecting to local Daemon
func createPact() dsl.Pact {
	pactDaemonPort := 6666
	return dsl.Pact{
		Port:     pactDaemonPort,
		Consumer: "MyConsumer",
		Provider: "MyProvider",
		LogDir:   logDir,
		PactDir:  pactDir,
	}
}

func TestPactConsumerLoginHandler_UserUnauthorised(t *testing.T) {

	var testBillyUnauthorized = func() error {
		err := AccessProvider(fmt.Sprintf("http://localhost:%d", pact.Server.Port))
		if err != nil {
			return fmt.Errorf("Error occured: %v", err)
		}
		return nil
	}

	pact.
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

	err := pact.Verify(testBillyUnauthorized)
	if err != nil {
		t.Fatalf("Error on Verify: %v", err)
	}

}
