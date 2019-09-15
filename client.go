package scpi

import (
	"context"
	"time"
)

// Client is a client of a device controlled using SCPI commands.
type Client interface {
	// Close closes the connection.
	Close() error

	// Exec executes a SCPI command.
	Exec(cmd string) error

	// ExecContext executes a SCPI command.
	ExecContext(ctx context.Context, cmd string) error

	// BulkExec executes multiple SCPI commands.
	BulkExec(cmds ...string) error

	// BulkExecContext executes multiple SCPI commands.
	BulkExecContext(ctx context.Context, cmds ...string) error

	// Ping verifies the connection to the device is still alive,
	// establishing a connection if necessary.
	Ping() error

	// PingContext verifies the connection to the device is still alive,
	// establishing a connection if necessary.
	PingContext(ctx context.Context) error

	// Query queries the device for the results of the specified command.
	Query(cmd string) (res string, err error)

	// QueryContext queries the device for the results of the specified command.
	QueryContext(ctx context.Context, cmd string) (res string, err error)
}

// NewClient returns a new client of a device controlled using SCPI commands.
func NewClient(proto, addr string, timeout time.Duration) (Client, error) {
	switch proto {
	case "tcp":
		return NewTCPClient(addr, timeout)
	default:
		return nil, InvalidProtocolError(proto)
	}
}
