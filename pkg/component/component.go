package component

import (
	"expvar"
	"fmt"
	"runtime"
	"sync"

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

	// 2) Setup NATS event callbacks
	//
	// - Error
	// - Reconnect
	// - Disconnected
	// - Closed
	// - Discovered
	//

	// 5) Register component so that it is available for discovery requests.
	//
	// - Subscribe to _NYFT.discovery
	//

	// 6) Register component so that it is available for direct status requests.
	//
	// - Subscribe to _NYFT.<component_id>.statsz
	//

	return nil
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

// Shutdown makes the component go away.
func (c *Component) Shutdown() error {
	c.NATS().Drain()
	return nil
}
