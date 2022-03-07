package listener

// OnEventFunc defines what to do when a `Listener` receives an event.
type OnEventFunc func(event interface{})

// IListener defines what a `Listener` does.
type IListener interface {
	// GetName returns the name of the `Listener`.
	GetName() string

	// SetChannel sets the communication channel between a `Broker` and a
	// `Listener`.
	SetChannel(chan interface{})

	// Listen and reacts to events.
	Listen()
}

// Listener reacts to events sent to its channel.
type Listener struct {
	// Name of the listener.
	Name string

	// OnEventFunc what to do when a `Listener` receives an event.
	OnEventFunc OnEventFunc

	// Ch is the channel which a `Broker` pushes events.
	Ch chan interface{}
}

// GetName returns the name of the listener.
func (s *Listener) GetName() string {
	return s.Name
}

// SetChannel sets the communication channel between a `Broker` and a
// `Listener`.
func (s *Listener) SetChannel(ch chan interface{}) {
	s.Ch = ch
}

// Listen and reacts to events.
func (s *Listener) Listen() {
	for {
		s.OnEventFunc(<-s.Ch)
	}
}

// New is the Listener factory.
func New(name string, onEventFunc OnEventFunc) *Listener {
	return &Listener{
		Name:        name,
		OnEventFunc: onEventFunc,
	}
}
