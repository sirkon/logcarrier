package main

import (
	"fmt"

	"github.com/sirkon/logcarrier/notify"
)

// NotifierType represents enumerated type of notifier
type NotifierType string

const (
	VoidNotifier NotifierType = ""
	FileNotifier NotifierType = "file"
)

var supportedNotifiers = map[NotifierType]func() notify.Notifier{
	VoidNotifier: func() notify.Notifier { return notify.VoidNotifier{} },
	FileNotifier: func() notify.Notifier { return notify.NewFileNotifier() },
}

// GetNotifier returns notifer of given type
func GetNotifier(nt NotifierType) (notify.Notifier, bool) {
	res, ok := supportedNotifiers[nt]
	return res(), ok
}

// UnmarshalText implementation
func (s *NotifierType) UnmarshalText(text []byte) error {
	t := string(text)
	_, ok := supportedNotifiers[NotifierType(t)]
	if !ok {
		return fmt.Errorf("Unsupported notifier type `\033[1m%s\033[0m`", t)
	}
	*s = NotifierType(t)
	return nil
}
