package niltime

import "time"

// Time implements json nil-able time.Time.
type Time struct {
	time.Time
}

// MarshalJSON overrides json marshaling rules for NilTime.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}

	return t.Time.MarshalJSON()
}
