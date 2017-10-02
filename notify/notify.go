package notify

// Notifier is the interface that wraps basic Notify method
type Notifier interface {
	// Notify sends fileName string to the a queue beneath
	Notify(fileName string) error

	// Init initializes notifier with configuration
	Init(data []byte) error
}
