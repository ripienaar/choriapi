package main

import (
	"encoding/json"
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

func NewDHT220Agent() (*DHT220Agent, error) {
	a := &DHT220Agent{
		log: logrus.WithFields(logrus.Fields{"agent": "dht220"}),
		meta: &agents.Metadata{
			Name:        "dht220",
			Description: "DHT220 Agent",
			Author:      "R.I.Pienaar <rip@devco.net>",
			Version:     "0.0.1",
			License:     "Apache-2.0",
			Timeout:     10,
			URL:         "http://choria.io",
		},
	}

	return a, nil
}

func (da *DHT220Agent) readingAction(req *RPCRequestBody, result *RPCReply) {
	var err error

	result.Data, err = rpi.read()
	if err != nil {
		result.Statuscode = 5
		result.Statusmsg = fmt.Sprintf("Could not read data to publish: %s", err.Error())
	}
}

// Everything below will go in some form of helper so does not need to be types by everyone

type RPCReply struct {
	Statuscode int         `json:"statuscode"`
	Statusmsg  string      `json:"statusmsg"`
	Data       interface{} `json:"data"`
}

type RPCRequestBody struct {
	Agent  string `json:"agent"`
	Action string `json:"action"`
}

func (da *DHT220Agent) Handle(msg *choria.Message, request protocol.Request, outbox chan *agents.AgentReply) {
	var err error

	rpcreply := da.newReply()
	defer da.publish(rpcreply, msg, request, outbox)

	rpcreq, err := da.requestFromMsg(msg.Payload)
	if err != nil {
		rpcreply.Statuscode = 5
		rpcreply.Statusmsg = fmt.Sprintf("Could not process request: %s", err.Error())

		return
	}

	switch rpcreq.Action {
	case "reading":
		da.readingAction(rpcreq, rpcreply)
	default:
		rpcreply.Statuscode = 2
		rpcreply.Statusmsg = fmt.Sprintf("Unknown action %s", rpcreq.Action)
	}
}

func (da *DHT220Agent) publish(rpc *RPCReply, msg *choria.Message, request protocol.Request, outbox chan *agents.AgentReply) {
	reply := &agents.AgentReply{
		Message: msg,
		Request: request,
	}

	j, err := json.Marshal(rpc)
	if err != nil {
		logrus.Errorf("Could not JSON encode reply: %s", err.Error())
		reply.Error = err
	}

	reply.Body = j

	outbox <- reply
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

func (da *DHT220Agent) requestFromMsg(msg string) (*RPCRequestBody, error) {
	r := &RPCRequestBody{}

	err := json.Unmarshal([]byte(msg), r)
	if err != nil {
		return nil, fmt.Errorf("Could not parse incoming request: %s", err.Error())
	}

	return r, nil
}
