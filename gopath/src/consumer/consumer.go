package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func main() {
	AccessProvider("http://127.0.0.1:8090")
}

func AccessProvider(providerUrl string) error {

	resp, err := http.Get(providerUrl + "/test.json")
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	fmt.Print(body)
	return nil

}
