package instrumentation

import (
    "net/http"
    "log"
)

func StartInstrumented(port string) {

    http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
        rw.Header().Add("Content-Type", "application/json; charset=utf-8")
        rw.Write([]byte(`{"foo":"bar"}`))
    })

    err := http.ListenAndServe(":" + port, nil)

    if err != nil {
        log.Fatal(err)
    }

}