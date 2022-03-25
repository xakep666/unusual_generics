package unusual_generics

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// SQLJSON is a wrapper for jsonb-like fields in sql databases.
// It marshals/unmarshalls internal type as json. See examples for usage.
type SQLJSON[T any] struct {
	X T // this way because we can't embed type parameter or use it as RHS
}

// SQLJSONOf constructs SQLJSON.
func SQLJSONOf[T any](v T) SQLJSON[T] {
	return SQLJSON[T]{X: v}
}

var (
	_ sql.Scanner   = &SQLJSON[any]{}
	_ driver.Valuer = SQLJSON[any]{}
)

// Scan implements sql.Scanner.
func (s *SQLJSON[T]) Scan(x any) error {
	switch t := x.(type) {
	case string:
		return json.Unmarshal([]byte(t), &s.X)
	case []byte:
		return json.Unmarshal(t, &s.X)
	default:
		return fmt.Errorf("unkown type %T", x)
	}
}

// Value implements driver.Valuer.
func (s SQLJSON[T]) Value() (driver.Value, error) {
	res, err := json.Marshal(s.X)
	if err != nil {
		return nil, err
	}

	return string(res), nil
}
