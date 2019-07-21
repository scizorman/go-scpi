package scpi

import (
	"context"
)

// Client is a client of a device controlled using SCPI commands.
type Client interface {
	// Close closes the connection.
	Close() error

	// Exec executes a SCPI command.
	Exec(cmd string) error

	// ExecContext executes a SCPI command.
	ExecContext(ctx context.Context, cmd string) error

	// Ping verifies the connection to the device is still alive,
	// establishing a connection if necessary.
	Ping() error

	// PingContext verifies the connection to the device is still alive,
	// establishing a connection if necessary.
	PingContext(ctx context.Context) error

	// Query queries the device for the results of the specified command.
	Query(cmd string) (res []byte, err error)

	// QueryContext queries the device for the results of the specified command.
	QueryContext(ctx context.Context, cmd string) (res []byte, err error)
}
