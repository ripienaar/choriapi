package main

import (
	"fmt"

	"github.com/choria-io/go-choria/server/agents/mcorpc"

	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-choria/server/agents"
	"github.com/sirupsen/logrus"
)

func NewDHT220Agent() (*mcorpc.Agent, error) {
	metadata := &agents.Metadata{
		Name:        "dht220",
		Description: "DHT220 Agent",
		Author:      "R.I.Pienaar <rip@devco.net>",
		Version:     "0.0.1",
		License:     "Apache-2.0",
		Timeout:     10,
		URL:         "http://choria.io",
	}

	agent := mcorpc.New("dht220", metadata, fw, logrus.WithFields(logrus.Fields{"agent": "dht220"}))
	err := agent.RegisterAction("reading", readingAction)

	return agent, err
}

func readingAction(req *mcorpc.Request, reply *mcorpc.Reply, agent *mcorpc.Agent, conn choria.ConnectorInfo) {
	var err error

	reply.Data, err = rpi.read()
	if err != nil {
		reply.Statuscode = 5
		reply.Statusmsg = fmt.Sprintf("Could not read data to publish: %s", err.Error())
	}
}
