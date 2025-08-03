package errorutil

import (
	"errors"
	"fmt"
)

func DeferErr(err *error, op func () error) {
	// nil error is not joined so no check is required
	*err = errors.Join(*err, op())
}

func DeferErrf(err *error, format string, op func () error) {
	if opErr := op(); opErr != nil {
		*err = errors.Join(*err, fmt.Errorf(format, opErr))
	}
}

func DeferIgnoreErr(op func () error) {
	_ = op()
}
