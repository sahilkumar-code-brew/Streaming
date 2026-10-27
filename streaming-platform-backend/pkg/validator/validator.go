package validator

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

var streamIDPattern = regexp.MustCompile(`^[a-z0-9-]{3,64}$`)

func ValidateStreamID(streamID string) error {
	if !streamIDPattern.MatchString(streamID) {
		return fmt.Errorf("stream_id must match ^[a-z0-9-]{3,64}$")
	}
	return nil
}

func ValidateDeviceID(deviceID string) error {
	parsed, err := uuid.Parse(deviceID)
	if err != nil || parsed.Version() != 4 {
		return fmt.Errorf("X-Device-Id must be a UUID v4")
	}
	return nil
}
