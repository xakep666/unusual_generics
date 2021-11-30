package unusual_generics

import (
	"time"
)

// TimeFormatted is a type that helps to deal with non-standard (RFC3339) time layouts in JSON/XML/etc.
// This type must be instantiated with interface implementation that returns layout and optionally location.
// Implementation must be able to work when it's zero-value.
// See 'time_format_test.go' for usage example.
type TimeFormatted[T TimeLayoutProvider] time.Time

// TimeLayoutProvider is a setup interface for TimeFormatted.
type TimeLayoutProvider interface {
	~struct{} // constrain base type to avoid runtime errors on zero-value method calls

	// TimeLayout must return time layout string in Go time format (see 'time' package).
	TimeLayout() string

	// TimeLocation can optionally return location that will be used for parsing and formatting.
	// This is useful when time format does not include timezone.
	TimeLocation() *time.Location
}

// FromTime constructs TimeFormatted from time.Time.
func FromTime[T TimeLayoutProvider](t time.Time) TimeFormatted[T] { return TimeFormatted[T](t) }

func (tf *TimeFormatted[T]) UnmarshalText(text []byte) error {
	var (
		err error
		t   time.Time
		p   T
	)

	if loc := p.TimeLocation(); loc != nil {
		t, err = time.ParseInLocation(p.TimeLayout(), string(text), loc)
	} else {
		t, err = time.Parse(p.TimeLayout(), string(text))
	}

	if err != nil {
		return err
	}

	*tf = TimeFormatted[T](t)

	return nil
}

func (tf TimeFormatted[T]) MarshalText() ([]byte, error) {
	var p T

	t := time.Time(tf)

	if loc := p.TimeLocation(); loc != nil {
		t = t.In(loc)
	}

	return []byte(t.Format(p.TimeLayout())), nil
}

// Time gives access to original time.Time to call methods on it.
func (tf TimeFormatted[T]) Time() time.Time { return time.Time(tf) }
