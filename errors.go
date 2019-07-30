package scpi

import (
	"fmt"
)

// InvalidProtocolError occures if the protocol is invalid.
type InvalidProtocolError string

func (e InvalidProtocolError) Error() string {
	return fmt.Sprintf("invalid protocol %s", e)
}
