// Copyright 2022 The broker Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package broker

import (
	"sync"

	"github.com/thalesfsp/broker/listener"
)

const subscribeChCapacity = 5

// Broker broadcast a message to all subscribed `Listeners`.
type Broker struct {
	// control flag.
	isRunning bool

	// Communication channel where events are broadcasted.
	publishCh chan interface{}

	// Stop the broker.
	stopCh chan struct{}

	// Subscriber channel registers a subscriber.
	subscriberCh chan chan interface{}

	// A list of subscribers.
	subscribers map[chan interface{}]struct{}

	// Unsubscriber channel unregister a subscriber.
	unsubscriberCh chan chan interface{}

	mu sync.Mutex
}

// GetIsRunning returns if the broker running, or not.
func (b *Broker) GetIsRunning() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.isRunning
}

// setIsRunning updates if the broker is running, or not.
func (b *Broker) setIsRunning(r bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.isRunning = r
}

// getStopCh returns the broker's stop channel. If closed, it stops the broker.
func (b *Broker) getStopCh() chan struct{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.stopCh
}

// setStopCh sets the broker's stop channel. It is used to set or reset the
// channel for when the broker is started or restarted.
func (b *Broker) setStopCh(ch chan struct{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.stopCh = ch
}

// Broker factory.
func New() *Broker {
	return &Broker{
		isRunning:      false,
		publishCh:      make(chan interface{}, 1),
		subscriberCh:   make(chan chan interface{}, 1),
		subscribers:    map[chan interface{}]struct{}{},
		unsubscriberCh: make(chan chan interface{}, 1),

		mu: sync.Mutex{},
	}
}

// Starts the broker.
//
// NOTE: It blocks execution. Run, or call from a goroutine.
func (b *Broker) Start() {
	//////
	// Prepares the broker to start
	//////
	b.setIsRunning(true)

	b.setStopCh(make(chan struct{}))

	// Should only run the loop if the broker should be running.
	for b.isRunning {
		select {
		// Stop the broker.
		case <-b.stopCh:
			b.setIsRunning(false)

			return

		// Add a subscriber.
		case subscriberCh := <-b.subscriberCh:
			b.subscribers[subscriberCh] = struct{}{}

		// Remove a subscriber.
		case unsubscriberCh := <-b.unsubscriberCh:
			delete(b.subscribers, unsubscriberCh)

		// Broadcast a message.
		case publishCh := <-b.publishCh:
			for sub := range b.subscribers {
				// NOTE: Buffered, use non-blocking send to protect the broker
				select {
				case sub <- publishCh:
				default:
				}
			}
		}
	}
}

// Stops the broker.
func (b *Broker) Stop() {
	close(b.getStopCh())
}

// Subscribe a `Listener`.
func (b *Broker) Subscribe(l listener.IListener) chan interface{} {
	ch := make(chan interface{}, subscribeChCapacity)

	// Subscriber <-> Listener bus.
	l.SetChannel(ch)

	go l.Listen()

	b.subscriberCh <- ch

	return ch
}

// Unsubscribe a `Listener`.
//
// NOTE: It blocks execution. Run, or call from a goroutine.
func (b *Broker) Unsubscribe(ch chan interface{}) {
	b.unsubscriberCh <- ch
}

// Publish a message.
func (b *Broker) Publish(msg interface{}) {
	if b.GetIsRunning() {
		b.publishCh <- msg
	}
}
