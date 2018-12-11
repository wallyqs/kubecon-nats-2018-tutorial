package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/wallyqs/kubecon-nats-2018-tutorial/pkg/component"
	"github.com/wallyqs/kubecon-nats-2018-tutorial/pkg/types"
)

const (
	Version = "0.1.0"
)

type Server struct {
	*component.Component
}

// SetupSubscriptions registers interest to the subjects that the
// Rides Manager will be handling.
func (s *Server) SetupSubscriptions() error {
	nc := s.NATS()

	// Helps finding an available driver to accept a drive request.
	nc.QueueSubscribe("drivers.find", "manager", func(msg *nats.Msg) {
		var req *types.DriverAgentRequest
		err := json.Unmarshal(msg.Data, &req)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
		log.Printf("requestID:%s - Driver Find Request\n", req.RequestID)
		response := &types.DriverAgentResponse{}

		// Find an available driver that can handle the user request.
		m, err := nc.Request("drivers.rides", msg.Data, 2*time.Second)
		if err != nil {
			response.Error = "No drivers available found, sorry!"
			resp, err := json.Marshal(response)
			if err != nil {
				log.Printf("requestID:%s - Error preparing response: %s",
					req.RequestID, err)
				return
			}

			// Reply with error response
			nc.Publish(msg.Reply, resp)
			return
		}
		response.ID = string(m.Data)

		resp, err := json.Marshal(response)
		if err != nil {
			response.Error = "No drivers available found, sorry!"
			resp, err := json.Marshal(response)
			if err != nil {
				log.Printf("requestID:%s - Error preparing response: %s",
					req.RequestID, err)
				return
			}

			// Reply with error response
			nc.Publish(msg.Reply, resp)
			return
		}
		log.Printf("requestID:%s - Driver Find Response: %+v\n",
			req.RequestID, string(m.Data))
		nc.Publish(msg.Reply, resp)
	})

	return nil
}
