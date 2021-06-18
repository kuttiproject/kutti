package cli

import "fmt"

// Error represents an error that can happen during kutti CLI
// execution. It can wrap a Go error or a simple message.
// It also has an exit code which will be returned to the OS
// if it bubbles to the top.
type Error struct {
	Err      error
	message  string
	Exitcode int
}

func (c *Error) Error() string {
	if c.Err != nil {
		return c.Err.Error()
	}

	return c.message
}

// WrapError wraps a Go error into a cli Error.
func WrapError(exitcode int, err error) error {
	return &Error{
		Err:      err,
		Exitcode: exitcode,
	}
}

// WrapErrorMessage wraps a string into a cli Error.
func WrapErrorMessage(exitcode int, message string) error {
	return &Error{
		message:  message,
		Exitcode: exitcode,
	}
}

// WrapErrorMessagef wraps a formatted string into a cli Error.
// Arguments are specified in the manner of fmt.Printf.
func WrapErrorMessagef(exitcode int, messageformat string, v ...interface{}) error {
	return &Error{
		message:  fmt.Sprintf(messageformat, v...),
		Exitcode: exitcode,
	}
}

// UnwrapError tries to typecast a Go error into a cli Error.
func UnwrapError(err error) (*Error, bool) {
	result, ok := err.(*Error)
	return result, ok
}
