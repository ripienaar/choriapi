package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/choria-io/go-choria/server/data"
)

// StartRegistration is the interface to registration in Choria
func (dh *RPi) StartRegistration(ctx context.Context, wg *sync.WaitGroup, interval int, output chan *data.RegistrationItem) {
	defer wg.Done()

	log.Printf("Starting to send data every %d seconds", interval)

	err := dh.publish(output)
	if err != nil {
		log.Printf("Could not create registration data: %s", err.Error())
	}

	for {
		select {
		case <-time.Tick(time.Duration(interval) * time.Second):
			err = dh.publish(output)
			if err != nil {
				log.Printf("Could not create registration data: %s", err.Error())
			}

		case <-ctx.Done():
			return
		}
	}
}

func (dh *RPi) publish(output chan *data.RegistrationItem) error {
	cur, err := dh.read()
	if err != nil {
		return err
	}

	j, err := json.Marshal(cur)
	if err != nil {
		return err
	}

	item := &data.RegistrationItem{
		Data:        &j,
		TargetAgent: "temperature",
	}

	log.Printf("Publishing %s", string(j))

	output <- item

	return nil
}
