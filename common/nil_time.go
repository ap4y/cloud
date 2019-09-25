package common

import "time"

// NilTime implements json nil-able time.Time.
type NilTime struct {
	time.Time
}

// MarshalJSON overrides json marshaling rules for NilTime.
func (t NilTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}

	return t.Time.MarshalJSON()
}
