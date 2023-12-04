package eerror

import (
	"fmt"
	"os"

	E "github.com/pkg/errors"
)

// Throw Print an error with the supplied message.
func Throw(msg string) {
	throw(msg, nil, 0)
}

// ThrowWrap Print an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
func ThrowWrap(msg string, err error) {
	throw(msg, err, 0)
}

// ThrowWrap Print an error message and
// then exit the program through the given code
func ThrowWithCode(msg string, code int) {
	throw(msg, nil, code)
}

func throw(msg string, err error, code int) {
	var errInfo error
	if err != nil {
		errInfo = E.Wrap(err, msg)
	} else {
		errInfo = E.New(msg)
	}

	fmt.Printf("Error: %+v \n", errInfo)
	os.Exit(code)
}
