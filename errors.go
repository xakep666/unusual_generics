package unusual_generics

import "errors"

type errorAs[T any] interface{ ErrorAs() T }

// ErrorAs is reflect(lite)-free analog of errors.As that not requires variable declaration for usage.
// It can be used in one-line 'if's. Second return value is true when Target type found in error chain, false otherwise.
// ErrorAs can also search types with method ErrorAs() Target
//  which implementation may be simple (without type assertions) in comparison with As(interface{}) bool.
func ErrorAs[Target comparable](err error) (Target, bool) {
	var ret Target

	for err != nil {
		if ret, ok := err.(Target); ok {
			return ret, true
		}

		if x, ok := err.(errorAs[Target]); ok {
			return x.ErrorAs(), true
		}

		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(&ret) {
			return ret, true
		}

		err = errors.Unwrap(err)
	}

	return ret, false
}

// ErrorAsPtr acts like ErrorAs but returns only one value. It's nil when Target type was not found in error chain.
func ErrorAsPtr[Target comparable](err error) *Target {
	var ret Target

	for err != nil {
		if ret, ok := err.(Target); ok {
			return &ret
		}

		if x, ok := err.(errorAs[Target]); ok {
			ret = x.ErrorAs()
			return &ret
		}

		if x, ok := err.(interface{ As(interface{}) bool }); ok && x.As(&ret) {
			return &ret
		}

		err = errors.Unwrap(err)
	}

	return nil
}

// ComparableError describes constraint for ErrorIs target.
type ComparableError interface {
	comparable
	error
}

type errorIs[T any] interface{ ErrorIs(T) bool }

// ErrorIs is analog of errors.Is implemented without reflect(lite).
// In comparison to errors.Is it can't be used like 'errors.Is(err, nil)' because 2nd argument can't be untyped nil.
// ErrorIs can also search types with method ErrorIs(Target) bool
//  which implementation may be simple in comparison with Is(error) bool.
func ErrorIs[Target ComparableError](err error, target Target) bool {
	for {
		if t, ok := err.(Target); ok && t == target {
			return true
		}

		if x, ok := err.(errorIs[Target]); ok && x.ErrorIs(target) {
			return true
		}

		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(target) {
			return true
		}

		if err = errors.Unwrap(err); err == nil {
			return false
		}
	}
}
