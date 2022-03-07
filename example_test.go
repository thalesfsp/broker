package broker_test

import (
	"fmt"
	"time"

	"github.com/thalesfsp/broker"
	"github.com/thalesfsp/broker/listener"
)

// Machine state, could be anything.
type Machine struct {
	// State of the machine.
	State string
}

// Starts everything.
func ExampleNew() {
	fmt.Println("Started...")

	//////
	// Broker
	//////

	b := broker.New()

	go b.Start()

	time.Sleep(1 * time.Second)

	//////
	// Listeners
	//////

	l1 := listener.New("l1", func(event interface{}) {
		tunnel, ok := event.(Machine)

		if ok {
			fmt.Printf("l1 Current tunnel status is %+v\n", tunnel.State)
		} else {
			fmt.Println("Not a tunnel, doing nothing")
		}
	})

	l2 := listener.New("l2", func(event interface{}) {
		tunnel, ok := event.(Machine)

		if ok {
			fmt.Printf("l2 Current tunnel status is %+v\n", tunnel.State)
		} else {
			fmt.Println("Not a tunnel, doing nothing")
		}
	})

	l3 := listener.New("l3", func(event interface{}) {
		tunnel, ok := event.(Machine)

		if ok {
			fmt.Printf("l3 Current tunnel status is %+v\n", tunnel.State)
		} else {
			fmt.Println("Not a tunnel, doing nothing")
		}
	})

	//////
	// Subscriber <-> Listeners.
	//////

	b.Subscribe(l1)

	b.Subscribe(l2)

	b.Subscribe(l3)

	time.Sleep(1 * time.Second)

	//////
	// Simulates publishing.
	//////

	// Start publishing messages:
	go func() {
		states := []string{"off", "on"}

		for _, state := range states {
			b.Publish(Machine{State: state})

			time.Sleep(10 * time.Millisecond)
		}
	}()

	time.Sleep(1 * time.Second)

	b.Stop()

	time.Sleep(1 * time.Second)

	// Start publishing messages:
	go func() {
		states := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

		for _, state := range states {
			b.Publish(Machine{State: state})

			time.Sleep(10 * time.Millisecond)
		}
	}()

	time.Sleep(1 * time.Second)

	go b.Start()

	time.Sleep(1 * time.Second)

	// Start publishing messages:
	go func() {
		states := []string{"ping", "pong"}

		for _, state := range states {
			b.Publish(Machine{State: state})

			time.Sleep(10 * time.Millisecond)
		}
	}()

	time.Sleep(1 * time.Second)

	b.Stop()

	time.Sleep(1 * time.Second)

	// output:
	// 1
	// l1 setup with success
	// l2 setup with success
	// l3 setup with success
	// Current tunnel status is off
	// Current tunnel status is off
	// Current tunnel status is off
	// Current tunnel status is on
	// Current tunnel status is on
	// Current tunnel status is on
}
