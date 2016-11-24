package main

import (
    "net/http"
    log "github.com/Sirupsen/logrus"
    "provider/instrumentation"
)

func main() {

    instrumentation.StartInstrumented("8080")

}

