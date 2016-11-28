package main

import (
    "provider/instrumentation"
)

func main() {
    instrumentation.StartInstrumented("8080")
}