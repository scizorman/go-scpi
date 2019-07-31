package scpi

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
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
		return newTCPClient(addr, timeout)
	default:
		return nil, InvalidProtocolError(proto)
	}
}

// TCPClient is an implementation of the Client interface for TCP network connections.
type TCPClient struct {
	conn *net.TCPConn
}

func newTCPClient(addr string, timeout time.Duration) (*TCPClient, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	d := net.Dialer{
		Timeout: timeout,
	}
	conn, err := d.Dial("tcp", tcpAddr.String())
	if err != nil {
		return nil, err
	}
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, fmt.Errorf("failed to cast %T to *net.TCPConn", conn)
	}
	client := &TCPClient{
		conn: tcpConn,
	}
	return client, nil
}

// Close implements the Client Close method.
func (c *TCPClient) Close() error {
	return c.conn.Close()
}

// Exec implements the Client Exec method.
func (c *TCPClient) Exec(cmd string) error {
	return c.ExecContext(context.Background(), cmd)
}

// ExecContext implements the Client ExecContext method.
func (c *TCPClient) ExecContext(ctx context.Context, cmd string) error {
	b := []byte(cmd + "\n")
	if _, err := c.conn.Write(b); err != nil {
		return err
	}
	return c.queryError(ctx, cmd)
}

var errorRegexp = regexp.MustCompile(`([+-]\d{1,3}),\"(.*?)\"`)

func (c *TCPClient) queryError(ctx context.Context, cmd string) error {
	res, err := c.Query("SYST:ERR?")
	if err != nil {
		return err
	}

	re := errorRegexp.Copy()
	g := re.FindStringSubmatch(res)
	if g == nil {
		return fmt.Errorf("invalid error format: %s", res)
	}

	code, err := strconv.Atoi(g[1])
	if err != nil {
		return err
	}
	msg := strings.ToLower(g[2])

	if code == 0 {
		return nil
	}
	return newCommandError(cmd, code, msg)
}

// BulkExec implements the Client BulkExec method.
func (c *TCPClient) BulkExec(cmds ...string) error {
	return c.BulkExecContext(context.Background(), cmds...)
}

// BulkExecContext implements the Client BulkExecContext method.
func (c *TCPClient) BulkExecContext(ctx context.Context, cmds ...string) error {
	cmd := strings.Join(cmds, ";")
	return c.ExecContext(ctx, cmd)
}

// Ping implements the Client Ping method.
func (c *TCPClient) Ping() error {
	return c.PingContext(context.Background())
}

// PingContext implements the Client PingContext method.
func (c *TCPClient) PingContext(ctx context.Context) error {
	// BUG(scizorman): PingContext is not implemented yet.
	return nil
}

// Query implements the Client Query method.
func (c *TCPClient) Query(cmd string) (res string, err error) {
	return c.QueryContext(context.Background(), cmd)
}

// QueryContext implements the Client QueryContext method.
func (c *TCPClient) QueryContext(ctx context.Context, cmd string) (res string, err error) {
	if err := c.ExecContext(ctx, cmd); err != nil {
		return "", err
	}

	buf := make([]byte, 1024)
	l, err := c.conn.Read(buf)
	if err != nil {
		return "", err
	}

	res = string(buf[:l])
	return res, nil
}
