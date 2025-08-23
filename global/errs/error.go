package errs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

type Error struct {
	Code   ErrorCode
	msg    string
	frame  xerrors.Frame
	err    error
	format string
	args   []interface{}
}

func (e *Error) Error() string {
	return fmt.Sprint(e)
}

func (e *Error) FormatError(p xerrors.Printer) (next error) {
	if e.msg == "" {
		p.Printf("Code: %v", e.Code)
	} else {
		p.Printf("%s", e.msg)
	}
	e.frame.Format(p)
	return e.err
}

func (e *Error) Format(s fmt.State, c rune) {
	xerrors.FormatError(e, s, c)
}

// Unwrap returns the error underlying the receiver, which may be nil.
func (e *Error) Unwrap() error {
	return e.err
}

func new(c ErrorCode, err error, callDepth int, msg string, format string, args []interface{}) *Error {
	return &Error{
		Code:   c,
		msg:    msg,
		frame:  xerrors.Caller(callDepth),
		err:    err,
		format: format,
		args:   args,
	}
}

// New returns a new error with the given code, underlying error and message. Pass 1
// for the call depth if New is called from the function raising the error; pass 2 if
// it is called from a helper function that was invoked by the original function; and
// so on.
func New(c ErrorCode, err error, callDepth int, msg string) *Error {
	return new(c, err, callDepth, msg, msg, make([]interface{}, 0))
}

// Newf uses format and args to format a message, then calls New.
func Newf(c ErrorCode, err error, format string, args ...any) *Error {
	return new(c, err, 2, fmt.Sprintf(format, args...), format, args)
}

// Wrapf detect the underlying error code, uses format and args to format a message, then calls New.
func Wrapf(err error, format string, args ...any) *Error {
	return new(Code(err), err, 2, fmt.Sprintf(format, args...), format, args)
}

func Code(err error) ErrorCode {
	if err == nil {
		return OK
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	if errors.Is(err, context.Canceled) {
		return Canceled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return DeadlineExceeded
	}
	return Unknown
}

func Message(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) {
		return e.msg
	}
	return ""
}

func Format(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) {
		return e.format
	}
	return ""
}

func Args(err error) []interface{} {
	if err == nil {
		return make([]interface{}, 0)
	}
	var e *Error
	if errors.As(err, &e) {
		return e.args
	}
	return make([]interface{}, 0)
}

type compositeErr struct {
	errs []error
}

func (c *compositeErr) Error() string {
	n := len(c.errs)
	if n == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i, e := range c.errs {
		sb.WriteString(e.Error())
		if i < n-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func Errors(errs []error) error {
	return &compositeErr{errs: errs}
}
