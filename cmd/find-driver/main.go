package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	payload := []byte(`{"type":"large"}`)
	c := &http.Client{Timeout: 2 * time.Second}
	resp, err := c.Post("http://127.0.0.1:9090/v1/rides", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[Response] ", string(body))
}
