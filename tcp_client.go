package scpi

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// TCPClient is an implementation of the Client interface for TCP network connections.
type TCPClient struct {
	conn *net.TCPConn
}

// NewTCPClient returns a new TCP client of a device controlled using SCPI commands.
func NewTCPClient(addr string, timeout time.Duration) (*TCPClient, error) {
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
	client := &TCPClient{
		conn: conn.(*net.TCPConn),
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
	if err := c.exec(ctx, cmd); err != nil {
		return fmt.Errorf("failed to execute the command '%s': %s", cmd, err)
	}
	return c.queryError(ctx, cmd)
}

func (c *TCPClient) exec(ctx context.Context, cmd string) error {
	b := []byte(cmd + "\n")
	if _, err := c.conn.Write(b); err != nil {
		return err
	}
	return nil
}

func (c *TCPClient) queryError(ctx context.Context, prevCmd string) error {
	res, err := c.QueryContext(ctx, "SYST:ERR?")
	if err != nil {
		return err
	}
	return confirmError(prevCmd, res)
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
	if err := c.exec(ctx, cmd); err != nil {
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

