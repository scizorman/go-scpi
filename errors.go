package scpi

import (
	"fmt"
)

// InvalidProtocolError occures if the protocol is invalid.
type InvalidProtocolError string

func (e InvalidProtocolError) Error() string {
	return fmt.Sprintf("invalid protocol %s", e)
}

// CommandError is the error of SCPI commands.
type CommandError struct {
	cmd  string
	code int
	msg  string
}

// Code returns the error code of a SCPI device.
func (e *CommandError) Code() int {
	return e.code
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("'%s' returned %d: %s", e.cmd, e.code, e.msg)
}

func newCommandError(cmd string, code int, msg string) *CommandError {
	return &CommandError{
		cmd:  cmd,
		code: code,
		msg:  msg,
	}
}
