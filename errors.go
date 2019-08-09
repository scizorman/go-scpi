package scpi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

var cmdErrRegexp = regexp.MustCompile(`([+-]\d+),\"(.*?)\"`)

func confirmError(cmd, errRes string) error {
	re := cmdErrRegexp.Copy()
	g := re.FindStringSubmatch(errRes)
	if g == nil {
		return fmt.Errorf("invalid error format: %s", errRes)
	}

	code, err := strconv.Atoi(g[1])
	if err != nil {
		return err
	}
	if code == 0 {
		return nil
	}

	msg := strings.ToLower(g[2])

	cmdErr := &CommandError{
		cmd:  cmd,
		code: code,
		msg:  msg,
	}
	return cmdErr
}
