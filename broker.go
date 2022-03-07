package broker

import (
	"fmt"
	"sync"

	"github.com/thalesfsp/broker/listener"
)

// Broker broadcast a message to all subscribed `Listeners`.
type Broker struct {
	mu sync.Mutex

	// Communication channel where events are broadcasted.
	publishCh chan interface{}

	// control flag.
	running bool

	// Stop the broker.
	stopCh chan struct{}

	// Subscriber channel registers a subscriber.
	subCh chan chan interface{}

	// A list of subscribers.
	subscribers map[chan interface{}]struct{}

	// Unsubscriber channel unregister a subscriber.
	unsubCh chan chan interface{}
}

func (b *Broker) GetStatus() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.running
}

func (b *Broker) SetStatus(r bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.running = r
}

func (b *Broker) getStopCh() chan struct{} {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.stopCh
}

func (b *Broker) SetStopCh(ch chan struct{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.stopCh = ch
}

// Broker factory.
func New() *Broker {
	return &Broker{
		mu:          sync.Mutex{},
		publishCh:   make(chan interface{}, 1),
		running:     false,
		subCh:       make(chan chan interface{}, 1),
		subscribers: map[chan interface{}]struct{}{},
		unsubCh:     make(chan chan interface{}, 1),
	}
}

// Starts the broker.
func (b *Broker) Start() {
	// Control flag.
	b.SetStatus(true)

	// Reset signal.
	b.SetStopCh(make(chan struct{}))

	fmt.Println("\nstarting the broker")
	fmt.Println("registered listeners channels", b.subscribers)

	for b.running {
		select {
		// Stop the broker by breaking the loop.
		case <-b.stopCh:
			b.SetStatus(false)

			fmt.Println("\nstopping the broker")
			fmt.Println("registered listeners channels", b.subscribers)

			return
		case subCh := <-b.subCh:
			fmt.Println("add new listener to the broker")

			b.subscribers[subCh] = struct{}{}
		case unsubCh := <-b.unsubCh:
			fmt.Println("removed listener from the broker")

			delete(b.subscribers, unsubCh)
		case publishCh := <-b.publishCh:
			for sub := range b.subscribers {
				fmt.Println("new event received, broadcasting...")

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
	msgCh := make(chan interface{}, 5)

	// Subscriber <-> Listener bus.
	l.SetChannel(msgCh)
	go l.Listen()

	b.subCh <- msgCh

	// TODO: Remove it.
	fmt.Println(l.GetName(), "setup with success")

	return msgCh
}

// Unsubscribe a `Listener`.
func (b *Broker) Unsubscribe(msgCh chan interface{}) {
	b.unsubCh <- msgCh
}

// Publish a message.
func (b *Broker) Publish(msg interface{}) {
	if b.GetStatus() {
		fmt.Println("publising to be broadcasted...")

		b.publishCh <- msg
	}
}
