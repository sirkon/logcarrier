package notify

// VoidNotifier dummy notifier
type VoidNotifier struct{}

// Notify ...
func (v VoidNotifier) Notify(fileName string) error {
	return nil
}

// Init ...
func (v VoidNotifier) Init(data []byte) error {
	return nil
}
