package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// TODO: Finish
	
	// Let the tool talk directly via http or NATS protocol.
	// 
	// The one using NATS would be able to talk to an imported service via NGS.
	// 
	c := &http.Client{Timeout: 2 * time.Second}
	resp, err := c.Post("https://127.0.0.1:9090/v1/drivers", []byte(""))
	if err != nil || resp == nil {
		log.Fatalf("Could not retrive geo location data: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	g := &geo{}
	if err := json.Unmarshal(body, &g); err != nil {
		log.Fatalf("Error unmarshalling geo: %v", err)
	}
	// resp, err := http.DefaultClient.Get("")
}
