package main

import (
	"github.com/pact-foundation/pact-go/dsl"
	"fmt"
	"log"
	"github.com/pact-foundation/pact-go/types"
	provider "provider/instrumentation"
	"testing"
)

func TestFoo(m *testing.T) {

	go provider.StartInstrumented("32000")

	pact := &dsl.Pact{
		Port:     6666,
		Consumer: "MyConsumer",
		Provider: "MyProvider",
	}
	defer pact.Teardown()

	err := pact.VerifyProvider(types.VerifyRequest{
		ProviderBaseURL:        "http://localhost:32000",
		BrokerURL:		"http://localhost",
	})

	if err != nil {
		log.Fatal("Error:", err)
	}

	fmt.Println("Test Passed!")

}