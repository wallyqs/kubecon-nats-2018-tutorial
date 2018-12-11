package component

import (
	"encoding/json"
	"expvar"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/nats-io/nuid"
)

// Component is contains reusable logic related to handling
// of the connection to NATS in the system.
type Component struct {
	// cmu is the lock from the component.
	cmu sync.Mutex

	// id is a unique identifier used for this component.
	id string

	// nc is the connection to NATS.
	nc *nats.Conn

	// kind is the type of component.
	kind string
}

// NewComponent creates a
func NewComponent(kind string) *Component {
	id := nuid.Next()
	return &Component{
		id:   id,
		kind: kind,
	}
}

type logger struct {
	pid int
	id  string
}

func (l logger) Write(b []byte) (int, error) {
	return fmt.Printf("[%d] %s - %s - %s", l.pid, time.Now().Format(time.StampMilli), l.id, string(b))
}

func (c *Component) SetupLogging() {
	log.SetFlags(0)
	l := &logger{
		pid: os.Getpid(),
		id:  c.ID(),
	}
	log.SetOutput(l)
}

// SetupConnectionToNATS connects to NATS and registers the event
// callbacks and makes it available for discovery requests as well.
func (c *Component) SetupConnectionToNATS(servers string, options ...nats.Option) error {
	// Label the connection with the kind and id from component.
	options = append(options, nats.Name(c.Name()))

	c.cmu.Lock()
	defer c.cmu.Unlock()

	// 1) Connect to NATS with customized options.
	//
	// - Connect
	//
	// Connect to NATS with customized options.
	nc, err := nats.Connect(servers, options...)
	if err != nil {
		return err
	}
	c.nc = nc

	// 2) Setup NATS event callbacks
	//
	// - Error
	// - Reconnect
	// - Disconnected
	// - Closed
	// - Discovered
	//
	nc.SetErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		log.Printf("NATS error: %s\n", err)
	})
	nc.SetReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected to NATS at %s", nc.ConnectedUrl())
	})
	nc.SetDisconnectHandler(func(_ *nats.Conn) {
		log.Println("Disconnected from NATS!")
	})
	nc.SetClosedHandler(func(_ *nats.Conn) {
		log.Println("NATS connection closed!")
		c.Exit()
	})

	// 5) Register component so that it is available for discovery requests.
	//
	// - Subscribe to _NYFT.discovery
	//
	_, err = c.nc.Subscribe("_NYFT.discovery", func(m *nats.Msg) {
		// Reply back directly with own name if requested.
		if m.Reply != "" {
			nc.Publish(m.Reply, []byte(c.ID()))
		} else {
			log.Println("[Discovery] No Reply inbox, skipping...")
		}
	})
	if err != nil {
		return err
	}

	// 6) Register component so that it is available for direct status requests.
	//
	// - Subscribe to _NYFT.<component_id>.statsz
	//
	statusSubject := fmt.Sprintf("_NYFT.%s.statz", c.id)
	_, err = c.nc.Subscribe(statusSubject, func(m *nats.Msg) {
		if m.Reply != "" {
			log.Println("[Status] Replying with status...")

			result, err := json.Marshal(c.Statsz())
			if err != nil {
				log.Printf("Error: %s\n", err)
				return
			}
			nc.Publish(m.Reply, result)
		} else {
			log.Println("[Status] No Reply inbox, skipping...")
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Component) SetupSignalHandlers() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for sig := range sigCh {
		log.Printf("Trapped '%v' signal", sig)

		switch sig {
		case syscall.SIGINT:
			// Flush NATS buffer then disconnect from NATS.
			c.Exit()
			return
		case syscall.SIGTERM:
			// Gracefully shutdown the component by
			// draining NATS connection.
			c.Shutdown()
			return
		}
	}
}

// Statsz are the latest stats from the component.
func (c *Component) Statsz() interface{} {
	// Add a couple of metrics from expvars.
	mem := expvar.Get("memstats").(expvar.Func)().(runtime.MemStats)
	cmdline := expvar.Get("cmdline").(expvar.Func)().([]string)
	return struct {
		Kind string   `json:"kind"`
		ID   string   `json:"id"`
		Cmd  []string `json:"cmdline"`
		Mem  uint64   `json:"mem"`
	}{
		Kind: c.kind,
		ID:   c.id,
		Cmd:  cmdline,
		Mem:  mem.HeapAlloc,
	}
}

// NATS returns the current NATS connection.
func (c *Component) NATS() *nats.Conn {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return c.nc
}

// ID returns the ID from the component.
func (c *Component) ID() string {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return c.id
}

// Name is the label used to identify the NATS connection.
func (c *Component) Name() string {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return fmt.Sprintf("%s:%s", c.kind, c.id)
}

// Shutdown makes the component go away gracefully.
func (c *Component) Shutdown() error {
	log.Println("Shutting down...")
	c.NATS().Drain()
	return nil
}

// Shutdown makes the component go away gracefully.
func (c *Component) Exit() error {
	log.Println("Exiting...")
	defer os.Exit(1)

	nc := c.NATS()
	if nc.IsClosed() {
		return nil
	}
	nc.Close()

	return nil
}
