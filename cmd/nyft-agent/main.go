package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/nats-io/go-nats"
	"github.com/wallyqs/kubecon-nats-2018-tutorial/pkg/component"
	"github.com/wallyqs/kubecon-nats-2018-tutorial/pkg/nyft-agent"
)

func main() {
	var (
		showHelp    bool
		showVersion bool
		natsServers string
		agentType   string
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: nyft-agent [options...]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Setup default flags
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.StringVar(&natsServers, "s", nats.DefaultURL, "List of NATS Servers to connect")
	flag.StringVar(&agentType, "type", "regular", "Kind of vehicle")
	flag.Parse()

	switch {
	case showHelp:
		flag.Usage()
		os.Exit(0)
	case showVersion:
		fmt.Fprintf(os.Stderr, "NYFT Driver Agent v%s\n", agent.Version)
		os.Exit(0)
	}

	// Register component
	comp := component.NewComponent("driver-agent")
	comp.SetupLogging()
	log.Printf("Starting NYFT Driver Agent version %s", agent.Version)

	// 3) Reconnection logic
	//
	// Set infinite retries to never stop reconnecting to an
	// available NATS server in case of an unreliable connection.
	//
	// Also set reconnection attempts to every 2 seconds.
	//
	options := []nats.Option{}
	err := comp.SetupConnectionToNATS(natsServers, options...)

	ag := agent.Agent{
		Component: comp,
		AgentType: agentType,
	}
	err = ag.SetupSubscriptions()
	if err != nil {
		log.Fatal(err)
	}

	runtime.Goexit()
}
