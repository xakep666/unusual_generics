package unusual_generics

// Ptr is a simple function that helps to get literal pointer in one-line.
func Ptr[T any](v T) *T { return &v }
