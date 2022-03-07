// Copyright 2022 The broker Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package broker_test

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/thalesfsp/broker"
	"github.com/thalesfsp/broker/listener"
)

//////
// A data structure to be sent as an event.
//////

// Machine state, could be anything.
type Machine struct {
	// State of the machine.
	State string
}

//////
// Golang's built-in example mechanism requires writing to output. The example
// will be using multiple goroutines. A "safe" buffer is needed.
//////

// Buffer is a goroutine safe `bytes.Buffer`.
type Buffer struct {
	buffer bytes.Buffer
	mutex  sync.Mutex
}

// Write appends the contents of p to the buffer, growing the buffer as needed. It returns
// the number of bytes written.
func (s *Buffer) Write(p []byte) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.buffer.Write(p)
}

// String returns the contents of the unread portion of the buffer
// as a string.  If the Buffer is a nil pointer, it returns "<nil>".
func (s *Buffer) String() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.buffer.String()
}

// Expected What's the `Expectation` for `Text` in the safe buffer.
type Expected struct {
	// Expectation if `text` is expected or not.
	Expectation bool

	// Text expected.
	Text string
}

// Contains check if `expected` is expected in the safe buffer.
func (s *Buffer) Contains(expected ...Expected) []bool {
	accumulator := []bool{}

	for _, expectation := range expected {
		accumulator = append(
			accumulator,
			strings.Contains(s.String(), expectation.Text) == expectation.Expectation,
		)
	}

	return accumulator
}

//////
// Example starts here.
//////

// Starts everything.
func ExampleNew() {
	//////
	// Broker setup, and start.
	//////

	b := broker.New()

	go b.Start()

	//////
	// Listeners setup, and start.
	//////

	// Safe buffer where output will be written.
	stringBuffer := &Buffer{
		buffer: bytes.Buffer{},
	}

	l1 := listener.New("l1", func(event interface{}) {
		tunnel, ok := event.(Machine)

		if ok {
			fmt.Fprintf(stringBuffer, "l1 -> tunnel status: %s\n", tunnel.State)
		}
	})

	l2 := listener.New("l2", func(event interface{}) {
		tunnel, ok := event.(Machine)

		if ok {
			fmt.Fprintf(stringBuffer, "l2 -> tunnel status: %s\n", tunnel.State)
		}
	})

	l3 := listener.New("l3", func(event interface{}) {
		tunnel, ok := event.(Machine)

		if ok {
			fmt.Fprintf(stringBuffer, "l3 -> tunnel status: %s\n", tunnel.State)
		}
	})

	//////
	// Subscriber <-> Listeners.
	//////

	b.Subscribe(l1)

	b.Subscribe(l2)

	b.Subscribe(l3)

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

	// Fake time.
	time.Sleep(500 * time.Millisecond)

	b.Stop()

	// Start publishing messages:
	go func() {
		states := []string{"off", "on"}

		for _, state := range states {
			b.Publish(Machine{State: state})

			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Fake time.
	time.Sleep(500 * time.Millisecond)

	//////
	// Simulates unsubscribing listeners.
	//////

	go b.Unsubscribe(l2.GetChannel())
	go b.Unsubscribe(l3.GetChannel())

	// Fake time.
	time.Sleep(500 * time.Millisecond)

	go b.Start()

	// Start publishing messages:
	go func() {
		states := []string{"ping", "pong"}

		for _, state := range states {
			b.Publish(Machine{State: state})

			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Fake time.
	time.Sleep(500 * time.Millisecond)

	b.Stop()

	fmt.Println(stringBuffer.Contains(
		Expected{Text: "l1 -> tunnel status: off", Expectation: true},
		Expected{Text: "l1 -> tunnel status: off", Expectation: true},
		Expected{Text: "l1 -> tunnel status: on", Expectation: true},
		Expected{Text: "l1 -> tunnel status: ping", Expectation: true},
		Expected{Text: "l1 -> tunnel status: pong", Expectation: true},
		Expected{Text: "l2 -> tunnel status: off", Expectation: true},
		Expected{Text: "l2 -> tunnel status: on", Expectation: true},
		Expected{Text: "l2 -> tunnel status: ping", Expectation: false},
		Expected{Text: "l2 -> tunnel status: pong", Expectation: false},
		Expected{Text: "l3 -> tunnel status: off", Expectation: true},
		Expected{Text: "l3 -> tunnel status: on", Expectation: true},
		Expected{Text: "l3 -> tunnel status: ping", Expectation: false},
		Expected{Text: "l3 -> tunnel status: pong", Expectation: false},
	))

	// output:
	// [true true true true true true true true true true true true true]
}
