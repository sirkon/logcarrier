package main

import (
	"fmt"
	"time"
)

// Duration is for time interval TOML unmarshaling
type Duration time.Duration

// UnmarshalText toml unmarshaling implementation
func (d *Duration) UnmarshalText(rawText []byte) error {
	res, err := time.ParseDuration(string(rawText))
	if err != nil {
		return fmt.Errorf("Cannot parse `%s` as duration", string(rawText))
	}
	*d = Duration(res)
	return nil
}
