package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/wallyqs/kubecon-nats-2018-tutorial/pkg/types"
)

const Version = "0.1.0"

func main() {
	var (
		showHelp     bool
		showVersion  bool
		natsServers  string
		httpEndpoint string
		vehicleType  string
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: find-driver [options...]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Setup default flags
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.StringVar(&natsServers, "s", "", "NATS Server")
	flag.StringVar(&natsServers, "nats", "", "NATS Server ")
	flag.StringVar(&httpEndpoint, "url", "http://127.0.0.1:9090/v1/rides", "HTTP Endpoint")
	flag.StringVar(&vehicleType, "type", "regular", "Vehicle type")
	flag.Parse()

	switch {
	case showHelp:
		flag.Usage()
		os.Exit(0)
	case showVersion:
		fmt.Fprintf(os.Stderr, "find-driver v%s\n", Version)
		os.Exit(0)
	}

	req := types.DriverAgentRequest{
		Type: vehicleType,
	}
	payload, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	var response []byte
	if natsServers != "" {
		nc, err := nats.Connect(natsServers)
		if err != nil {
			log.Fatal(err)
		}
		msg, err := nc.Request("drivers.find", payload, 2*time.Second)
		if err != nil {
			log.Fatal(err)
		}
		response = msg.Data
	} else {
		// Use http by default
		c := &http.Client{Timeout: 2 * time.Second}
		resp, err := c.Post(httpEndpoint, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		response, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("[Response] ", string(response))
}
