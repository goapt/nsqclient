package nsqclient

import (
	"fmt"
)

type debugError struct {
	s string
}

func NewDebugError(msg interface{}) error {
	var message string
	if m, ok := msg.(error); ok {
		message = m.Error()
	} else {
		message = fmt.Sprint(msg)
	}

	return &debugError{message}
}

func (e *debugError) Error() string {
	return e.s
}
