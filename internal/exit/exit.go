package exit

import (
	"fmt"
	"os"
)

type Status uint8

const (
	StatusOk Status = 0
	// Unknown errors, service should be restarted
	// https://tldp.org/LDP/abs/html/exitcodes.html
	StatusFatal Status = 1
	StatusPanic Status = 2
	// Known errors, restart won’t be successful
	StatusGeneralError  Status = 11
	StatusInvalidFlag   Status = 12
	StatusInvalidArg    Status = 13
	StatusInvalidConfig Status = 14
	// ... ... ...
)

func Exit(status Status, errors ...error) {
	for _, e := range errors {
		fmt.Fprintln(os.Stderr, "error: "+e.Error())
	}
	os.Exit(int(status))
}
