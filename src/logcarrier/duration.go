package main

import (
	"fmt"
	"time"
)

// Duration is for time interval TOML unmarshaling
type Duration time.Duration

// UnmarshalYAML implementation
func (d *Duration) UnmarshalYAML(data []byte) error {
	res, err := time.ParseDuration(string(data))
	if err != nil {
		return fmt.Errorf("Cannot parse `%s` as duration", string(data))
	}
	*d = Duration(res)
	return nil

}
