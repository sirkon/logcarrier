package main

import (
	"fmt"
)

// LogrotateMethod describes methods of log rotation
type LogrotateMethod int

const (
	// LogrotatePeriodic means to rely on periodic automatic log
	// rotation
	LogrotatePeriodic LogrotateMethod = iota

	// LogrotateGuided means to rely on guided log rotation via
	// the carrier protoocol
	LogrotateGuided

	// LogrotateBoth means both periodic and guided log rotations
	// methods are allowed
	LogrotateBoth
)

func (lm LogrotateMethod) String() string {
	switch lm {
	case LogrotatePeriodic:
		return "periodic"
	case LogrotateGuided:
		return "guided"
	case LogrotateBoth:
		return "both"
	default:
		panic(fmt.Errorf("unsupported log rotation method %d", lm))
	}
}

var rotMethodMap = map[string]LogrotateMethod{
	LogrotatePeriodic.String(): LogrotatePeriodic,
	LogrotateGuided.String():   LogrotateGuided,
	LogrotateBoth.String():     LogrotateBoth,
}

// UnmarshalText yaml unmarshalling implementation
func (lm *LogrotateMethod) UnmarshalText(rawtext []byte) error {
	res, ok := rotMethodMap[string(rawtext)]
	if !ok {
		return fmt.Errorf("unsupported rotation method `\033[1m%s\033[0m`", string(rawtext))
	}
	*lm = res
	return nil
}
