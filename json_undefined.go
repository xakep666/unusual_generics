package unusual_generics

import "encoding/json"

// JSONUndefined is a type that helps to distinguish if field was explicitly defined in JSON or not.
// JavaScript has special keyword 'undefined' to handle such cases.
// See examples for better explanation.
type JSONUndefined[T any] struct {
	// Defined set to true if value was explicitly defined in JSON.
	Defined bool

	Value *T
}

func (j *JSONUndefined[T]) UnmarshalJSON(bytes []byte) error {
	j.Defined = true // UnmarshalJSON will not be called if value was not explicitly defined.

	return json.Unmarshal(bytes, &j.Value)
}
