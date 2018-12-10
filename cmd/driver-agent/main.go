package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/nats-io/go-nats"
	"github.com/wallyqs/kubecon-nats-2018-tutorial/pkg/component"
	"github.com/wallyqs/kubecon-nats-2018-tutorial/pkg/driver-agent"
)

func main() {
	var (
		showHelp    bool
		showVersion bool
		natsServers string
		agentType   string
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: driver-agent [options...]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Setup default flags
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.StringVar(&natsServers, "nats", nats.DefaultURL, "List of NATS Servers to connect")
	flag.StringVar(&agentType, "type", "regular", "Kind of vehicle")
	flag.Parse()

	switch {
	case showHelp:
		flag.Usage()
		os.Exit(0)
	case showVersion:
		fmt.Fprintf(os.Stderr, "NATS Rider Driver Agent v%s\n", driveragent.Version)
		os.Exit(0)
	}
	log.Printf("Starting NATS Rider Driver Agent version %s", driveragent.Version)

	comp := component.NewComponent("driver-agent")

	// Set infinite retries to never stop reconnecting to an
	// available NATS server in case of an unreliable connection.
	err := comp.SetupConnectionToNATS(natsServers, nats.MaxReconnects(-1))
	if err != nil {
		log.Fatal(err)
	}

	ag := driveragent.Agent{
		Component: comp,
		AgentType: agentType,
	}
	err = ag.SetupSubscriptions()
	if err != nil {
		log.Fatal(err)
	}

	runtime.Goexit()
}
