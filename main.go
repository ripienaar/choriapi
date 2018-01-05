package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/choria-io/go-choria/build"
	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-choria/server"
	"github.com/choria-io/go-choria/server/data"
	"github.com/morus12/dht22"
)

type RPi struct {
	Pin string

	sensor         *dht22.DHT22
	choriaInstance *server.Instance
}

type reading struct {
	Temperature float32   `json:"temperature"`
	Humidity    float32   `json:"humidity"`
	Time        time.Time `json:"time"`
}

var rpi *RPi
var fw *choria.Framework
var mu *sync.Mutex

// NewRPi sets up the DHT22 reading and configures the embedded Choria
func NewRPi(pin string) (*RPi, error) {
	rpi := &RPi{
		Pin:    pin,
		sensor: dht22.New(pin),
	}

	mu = &sync.Mutex{}

	cfg, err := choria.NewConfig("/dev/null")
	if err != nil {
		return nil, err
	}

	cfg.DisableTLS = true
	build.Secure = "false"
	cfg.Choria.MiddlewareHosts = []string{"demo.nats.io:4222"}
	cfg.RegisterInterval = 60
	cfg.LogLevel = "debug"

	fw, err := choria.NewWithConfig(cfg)
	if err != nil {
		return nil, err
	}

	fw.SetupLogging(true)

	rpi.choriaInstance, err = server.NewInstance(fw)
	if err != nil {
		return nil, err
	}

	return rpi, nil
}

// Run starts the Choria instance and once connected register
// our agent and registration provider
func (dh *RPi) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	dh.choriaInstance.Run(ctx, wg)

	agent, err := NewDHT220Agent()
	if err != nil {
		log.Printf("Could not register DHT220 Agent: %s", err)
		panic(err)
	}

	err = dh.choriaInstance.RegisterAgent(ctx, "dht220", agent)
	if err != nil {
		log.Printf("Could not register DHT220 Agent: %s", err)
		panic(err)
	}

	dh.choriaInstance.RegisterRegistrationProvider(ctx, wg, dh)
}

func (dh *RPi) read() (*reading, error) {
	mu.Lock()
	defer mu.Unlock()

	temp, err := dh.sensor.Temperature()
	if err != nil {
		return nil, err
	}

	humidity, err := dh.sensor.Humidity()
	if err != nil {
		return nil, err
	}

	r := reading{
		Humidity:    humidity,
		Temperature: temp,
		Time:        time.Now(),
	}

	return &r, nil
}

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

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var err error

	rpi, err = NewRPi(os.Getenv("DH2200_PIN"))
	if err != nil {
		fmt.Printf(err.Error())
		cancel()
		return
	}

	rpi.Run(ctx, wg)

	for {
		select {
		case sig := <-sigs:
			log.Printf("Shutting down on %s", sig)
			cancel()
		case <-ctx.Done():
			wg.Wait()
			cancel()
			return
		}
	}
}
