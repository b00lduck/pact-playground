package main

import (
    "github.com/pact-foundation/pact-go/dsl"
	"fmt"
	"log"
    "github.com/pact-foundation/pact-go/types"
    provider "provider/instrumentation"
)

func main() {

    go provider.StartInstrumented("32000")

    pact := &dsl.Pact{
        Port:     6666,
        Consumer: "MyConsumer",
        Provider: "MyProvider",

    }
    defer pact.Teardown()

    err := pact.VerifyProvider(types.VerifyRequest{
        ProviderBaseURL:        "http://localhost:32000",
        PactURLs:               []string{"/home/daniel/dev/ws_gobama/pact-playground/gopath/src/provider_pact_test/pacts/myconsumer-myprovider.json"},

    })

    if err != nil {
        log.Fatal("Error:", err)
    }

    fmt.Println("Test Passed!")

}