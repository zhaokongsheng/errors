// Package errors provides simple error handling primitives.
//
// The traditional error handling idiom in Go is roughly akin to
//
//     if err != nil {
//             return err
//     }
//
// which when applied recursively up the call stack results in error reports
// without context or debugging information. The errors package allows
// programmers to add context to the failure path in their code in a way
// that does not destroy the original value of the error.
//
// Adding context to an error
//
// The errors.Wrap function returns a new error that adds context to the
// original error by recording a stack trace at the point Wrap is called,
// together with the supplied message. For example
//
//     _, err := ioutil.ReadAll(r)
//     if err != nil {
//             return errors.Wrap(err, "read failed")
//     }
//
// Formatted printing of errors
//
// All error values returned from this package implement fmt.Formatter and can
// be formatted by the fmt package. The following verbs are supported:
//
//     %s    print the error. If the error has a wrappedError it will be
//           printed recursively.
//     %v    see %s
//     %+v   extended format. Each Frame of the error's StackTrace will
//           be printed in detail.
//
// Retrieving the stack trace of an error or wrapper
//
// New, Errorf, Wrap, and Wrapf record a stack trace at the point they are
// invoked. This information can be retrieved with the following interface:
//
//     type stackTracer interface {
//             StackTrace() errors.StackTrace
//     }
//
// The returned errors.StackTrace type is defined as
//
//     type StackTrace []Frame
//
// The Frame type represents a call site in the stack trace. Frame supports
// the fmt.Formatter interface that can be used for printing information about
// the stack trace of this error. For example:
//
//     if err, ok := err.(stackTracer); ok {
//             for _, f := range err.StackTrace() {
//                     fmt.Printf("%+s:%d\n", f, f)
//             }
//     }
//
// Although the stackTracer interface is not exported by this package, it is
// considered a part of its stable public interface.
//
// See the documentation for Frame.Format for more details.
package errors

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func New(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// fundamental is an error that has a message and a stack, but no caller.
type fundamental struct {
	msg string
	*stack
}

func (f *fundamental) Error() string { return f.msg }

func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

type wrappedError struct {
	err   error
	msg   string
	stack *stack
}

func (w *wrappedError) Error() string { return w.msg + ": " + w.err.Error() }

func (w *wrappedError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, w.msg)

			isRoot := (errors.Unwrap(w.err) == nil)
			_, isFormatter := w.err.(fmt.Formatter)
			// Print the wrapped error message between the wrapping message and stack trace
			// if the wrapped error is the root error and does not implement fmt.Formatter.
			if isRoot && !isFormatter {
				fmt.Fprintf(s, "\n%+v", w.err)
				w.stack.Format(s, verb)
				io.WriteString(s, "\n")
			} else {
				w.stack.Format(s, verb)
				io.WriteString(s, "\n")
				fmt.Fprintf(s, "%+v\n", w.err)
			}

			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

func (w *wrappedError) Unwrap() error {
	return w.err
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied messages.
// The messages are join into one message with "\n" separator.
// If err is nil, Wrap returns nil.
func Wrap(err error, messages ...string) error {
	if err == nil {
		return nil
	}

	message := strings.Join(messages, "\n")
	_, ok := err.(fmt.Formatter)
	// If err already implements fmt.Formatter, add only the top stack trace
	if ok {
		return &wrappedError{
			err:   err,
			msg:   message,
			stack: topCaller(),
		}
	}

	return &wrappedError{
		err:   err,
		msg:   message,
		stack: callers(),
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	message := fmt.Sprintf(format, args...)

	_, ok := err.(fmt.Formatter)
	// If err already implements fmt.Formatter, add only the top stack trace
	if ok {
		return &wrappedError{
			err:   err,
			msg:   message,
			stack: topCaller(),
		}
	}

	return &wrappedError{
		err:   err,
		msg:   message,
		stack: callers(),
	}
}
