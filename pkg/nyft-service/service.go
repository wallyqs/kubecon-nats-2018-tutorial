package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/nats-io/nuid"
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
	// nc := s.NATS()

	// TUTORIAL) Use a load balanced QueueSubscription on 'drivers.find'
	// using the 'manager' group.

	return nil
}

func (s *Server) processFindRequest(msg *nats.Msg) {
	nc := s.NATS()

	var req *types.DriverAgentRequest
	err := json.Unmarshal(msg.Data, &req)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	// If no request ID, generate one of the request.
	if req.RequestID == "" {
		req.RequestID = nuid.Next()
	}

	log.Printf("requestID:%s - Driver Find Request\n", req.RequestID)
	response := &types.DriverAgentResponse{}

	// TUTORIAL) Find an available driver that can handle the user request.
	// 
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

	// TUTORIAL) Send the response back to requestor.
	//
	nc.Publish(msg.Reply, resp)
}
