package unusual_generics

import (
	"time"
)

// TimeFormatted is a type that helps to deal with non-standard (RFC3339) time layouts in JSON/XML/etc.
// This type must be instantiated with interface implementation that returns layout and optionally location.
// Implementation must be able to work when it's zero-value.
// See 'time_format_test.go' for usage example.
type TimeFormatted[T TimeLayoutProvider] struct {
	Time time.Time

	layoutProvider T
}

// TimeLayoutProvider is a setup interface for TimeFormatted.
type TimeLayoutProvider interface {
	// TimeLayout must return time layout string in Go time format (see 'time' package).
	TimeLayout() string

	// TimeLocation can optionally return location that will be used for parsing and formatting.
	// This is useful when time format does not include timezone.
	TimeLocation() *time.Location
}

func (tf *TimeFormatted[T]) UnmarshalText(text []byte) error {
	var err error

	if loc := tf.layoutProvider.TimeLocation(); loc != nil {
		tf.Time, err = time.ParseInLocation(tf.layoutProvider.TimeLayout(), string(text), loc)
	} else {
		tf.Time, err = time.Parse(tf.layoutProvider.TimeLayout(), string(text))
	}

	return err
}

func (tf TimeFormatted[T]) MarshalText() ([]byte, error) {
	t := tf.Time
	if loc := tf.layoutProvider.TimeLocation(); loc != nil {
		t = t.In(loc)
	}

	return []byte(t.Format(tf.layoutProvider.TimeLayout())), nil
}
