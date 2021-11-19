package unusual_generics_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/xakep666/unusual_generics"
)

type testErrorIsCase[Target unusual_generics.ComparableError] struct {
	err    error
	match  bool
	target Target
}

func testErrorIs[Target unusual_generics.ComparableError](t *testing.T, tc testErrorIsCase[Target]) {
	t.Helper()

	t.Run("", func(t *testing.T) {
		if got := unusual_generics.ErrorIs(tc.err, tc.target); got != tc.match {
			t.Errorf("ErrorIs(%v, %v) = %v, want %v", tc.err, tc.target, got, tc.match)
		}
	})
}

func TestErrorIs(t *testing.T) {
	err1 := errors.New("1")
	erra := wrapped{"wrap 2", err1}
	errb := wrapped{"wrap 3", erra}

	err3 := errors.New("3")

	poser := &poser{"either 1 or 3", func(err error) bool {
		return err == err1 || err == err3
	}}

	testErrorIs(t, testErrorIsCase[error]{err1, true, err1})
	testErrorIs(t, testErrorIsCase[error]{erra, true, err1})
	testErrorIs(t, testErrorIsCase[error]{errb, true, err1})
	testErrorIs(t, testErrorIsCase[error]{err1, false, err3})
	testErrorIs(t, testErrorIsCase[error]{erra, false, err3})
	testErrorIs(t, testErrorIsCase[error]{errb, false, err3})
	testErrorIs(t, testErrorIsCase[error]{poser, true, err1})
	testErrorIs(t, testErrorIsCase[error]{poser, true, err3})
	testErrorIs(t, testErrorIsCase[error]{poser, false, erra})
	testErrorIs(t, testErrorIsCase[error]{poser, false, errb})
	testErrorIs(t, testErrorIsCase[error]{errorUncomparable{}, false, err1})
	testErrorIs(t, testErrorIsCase[error]{&errorUncomparable{}, false, err1})
	testErrorIs(t, testErrorIsCase[*fs.PathError]{
		errIsOther{&fs.PathError{Op: "test"}},
		true,
		&fs.PathError{Op: "test"},
	})
	testErrorIs(t, testErrorIsCase[*fs.PathError]{
		errIsOther{&fs.PathError{Op: "test"}},
		false,
		&fs.PathError{Op: "test2"},
	})
}

type poser struct {
	msg string
	f   func(error) bool
}

var poserPathErr = &fs.PathError{Op: "poser"}

func (p *poser) Error() string     { return p.msg }
func (p *poser) Is(err error) bool { return p.f(err) }
func (p *poser) As(err interface{}) bool {
	switch x := err.(type) {
	case **poser:
		*x = p
	case *errorT:
		*x = errorT{"poser"}
	case **fs.PathError:
		*x = poserPathErr
	default:
		return false
	}
	return true
}

type errAsOther struct{}

func (e errAsOther) Error() string { return "errAsOther" }

func (e errAsOther) ErrorAs() *fs.PathError {
	return poserPathErr
}

type errIsOther struct {
	other *fs.PathError
}

func (e errIsOther) Error() string { return "errIsOther" }

func (e errIsOther) ErrorIs(pe *fs.PathError) bool { return pe.Op == e.other.Op }

type testErrorAsCase[Target comparable] struct {
	err   error
	match bool
	want  Target
}

func testErrorAs[Target comparable](t *testing.T, tc testErrorAsCase[Target]) {
	// dedicated function to use generics
	t.Helper()

	name := fmt.Sprintf("ErrorAs[%T](Errorf(..., %v))", tc.want, tc.err)
	t.Run(name, func(t *testing.T) {
		got, match := unusual_generics.ErrorAs[Target](tc.err)
		if match != tc.match {
			t.Fatalf("match: got %v; want %v", match, tc.match)
		}
		if !match {
			return
		}

		if got != tc.want {
			t.Fatalf("got %#v, want %#v", got, tc.want)
		}
	})

	name = fmt.Sprintf("ErrorAsPtr[%T](Errorf(..., %v))", tc.want, tc.err)
	t.Run(name, func(t *testing.T) {
		got := unusual_generics.ErrorAsPtr[Target](tc.err)
		if tc.match && got == nil {
			t.Fatalf("match: got nil; want non-nil")
		}
		if got == nil {
			return
		}

		if *got != tc.want {
			t.Fatalf("got %#v, want %#v", got, tc.want)
		}
	})
}

type timeout interface {
	Timeout() bool
}

func TestErrorAs(t *testing.T) {
	_, errF := os.Open("non-existing")
	poserErr := &poser{"oh no", nil}

	testErrorAs(t, testErrorAsCase[*fs.PathError]{
		err:   nil,
		match: false,
	})

	testErrorAs(t, testErrorAsCase[errorT]{
		err:   wrapped{"pitted the fool", errorT{"T"}},
		match: true,
		want:  errorT{"T"},
	})

	testErrorAs(t, testErrorAsCase[*fs.PathError]{
		err:   errF,
		match: true,
		want:  errF.(*fs.PathError),
	})

	testErrorAs(t, testErrorAsCase[*fs.PathError]{
		err:   errorT{},
		match: false,
	})

	testErrorAs(t, testErrorAsCase[errorT]{
		err:   wrapped{"wrapped", nil},
		match: false,
	})

	testErrorAs(t, testErrorAsCase[errorT]{
		err:   &poser{"error", nil},
		match: true,
		want:  errorT{"poser"},
	})

	testErrorAs(t, testErrorAsCase[*fs.PathError]{
		err:   &poser{"path", nil},
		match: true,
		want:  poserPathErr,
	})

	testErrorAs(t, testErrorAsCase[*poser]{
		err:   poserErr,
		match: true,
		want:  poserErr,
	})

	testErrorAs(t, testErrorAsCase[timeout]{
		err:   errors.New("err"),
		match: false,
	})

	testErrorAs(t, testErrorAsCase[timeout]{
		err:   errF,
		match: true,
		want:  errF.(timeout),
	})

	testErrorAs(t, testErrorAsCase[timeout]{
		err:   wrapped{"path error", errF},
		match: true,
		want:  errF.(timeout),
	})

	testErrorAs(t, testErrorAsCase[*fs.PathError]{
		err:   errAsOther{},
		match: true,
		want:  poserPathErr,
	})

	testErrorAs(t, testErrorAsCase[timeout]{
		err:   errAsOther{},
		match: false,
	})
}

type errorT struct{ s string }

func (e errorT) Error() string { return fmt.Sprintf("errorT(%s)", e.s) }

type wrapped struct {
	msg string
	err error
}

func (e wrapped) Error() string { return e.msg }

func (e wrapped) Unwrap() error { return e.err }

type errorUncomparable struct {
	f []string
}

func (errorUncomparable) Error() string {
	return "uncomparable error"
}

func (errorUncomparable) Is(target error) bool {
	_, ok := target.(errorUncomparable)
	return ok
}
