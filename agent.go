package main

import (
	"fmt"

	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-choria/protocol"
	"github.com/choria-io/go-choria/server/agents"
	"github.com/sirupsen/logrus"
)

type DHT220Agent struct {
	meta *agents.Metadata
	log  *logrus.Entry
}

type RPCRequestBody struct {
	Agent  string `json:"agent"`
	Action string `json:"action"`
}

type RPCReply struct {
	Statuscode int      `json:"statuscode"`
	Statusmsg  string   `json:"statusmsg"`
	Data       *reading `json:"data"`
}

func NewDHT220Agent() (*DHT220Agent, error) {
	a := &DHT220Agent{
		log: logrus.WithFields(logrus.Fields{"agent": "dht220"}),
		meta: &agents.Metadata{
			Name:        "dht220",
			Description: "DHT220 Agent",
			Author:      "R.I.Pienaar <rip@devco.net>",
			Version:     "0.0.1",
			License:     "Apache-2.0",
			Timeout:     2,
			URL:         "http://choria.io",
		},
	}

	return a, nil
}

func (da *DHT220Agent) Name() string {
	return da.meta.Name
}

func (da *DHT220Agent) Metadata() *agents.Metadata {
	return da.meta
}

func (da *DHT220Agent) newReply() *RPCReply {
	reply := &RPCReply{
		Statuscode: 0,
		Statusmsg:  "OK",
		Data:       &reading{},
	}

	return reply
}
func (da *DHT220Agent) Handle(msg *choria.Message, request protocol.Request, result chan *agents.AgentReply) {
	reply := &agents.AgentReply{
		Message: msg,
		Request: request,
	}

	var err error

	rpc := da.newReply()
	rpc.Data, err = rpi.read()
	if err != nil {
		rpc.Statuscode = 5
		rpc.Statusmsg = fmt.Sprintf("Could not read data to publish: %s", err.Error())
		reply.Error = fmt.Errorf(rpc.Statusmsg)
	}

	result <- reply
}
