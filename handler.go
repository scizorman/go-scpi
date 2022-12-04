package scpi

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Handler is a handler for a device controlled using SCPI commands.
type Handler struct {
	Client
}

// NewHandler returns a new handler for a device controlled using SCPI commands.
func NewHandler(client Client) *Handler {
	return &Handler{
		Client: client,
	}
}

// Reset resets the instrument to a factory pre-defined condition and clears the error log.
func (h *Handler) Reset() error {
	return h.BulkExec("*RST", "CLS")
}

// WaitForComplete waits for all queued operations to complete up to the specified timeout.
func (h *Handler) WaitForComplete(ctx context.Context, timeout time.Duration) error {
	subCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ch := make(chan error, 1)
	go func() {
		_, err := h.QueryContext(subCtx, "*WAI;*OPC?")
		ch <- err
	}()

	select {
	case err := <-ch:
		return err
	case <-subCtx.Done():
		// TODO(scizorman): Refactor the error
		return subCtx.Err()
	}
}

// Trigger triggers the device if, and only if,
// Bus Triggering is the type of trigger event selected.
// Otherwise, this command is ignored.
func (h *Handler) Trigger(ctx context.Context) error {
	return h.ExecContext(ctx, "*TRG")
}

// Identify returns the identification data.
//
// The standards order is follows:
//
//	Manufacturer
//	Model number
//	Serial number (or 0)
//	Firmware version
func (h *Handler) Identify(ctx context.Context) (id string, err error) {
	res, err := h.QueryContext(ctx, "*IDN?")
	if err != nil {
		return "", nil
	}

	id = string(res)
	return id, nil
}

// SetEventStatusEnable sets the value in the enable register for the Standard Event Status group.
// The selected bits are then reported to bit 5 of the Status Byte.
func (h *Handler) SetEventStatusEnable(ctx context.Context, bit uint8) error {
	cmd := fmt.Sprintf("*ESE %d", bit)
	return h.Exec(cmd)
}

// QueryEventStatusEnable queries the event status enable.
func (h *Handler) QueryEventStatusEnable(ctx context.Context) (bit uint8, err error) {
	res, err := h.QueryContext(ctx, "*ESE?")
	if err != nil {
		return 0, err
	}

	return parseBit(res)
}

// QueryEventStatusRegister queries the event status register.
// The register is cleared when it is executed.
func (h *Handler) QueryEventStatusRegister(ctx context.Context) (bit uint8, err error) {
	res, err := h.QueryContext(ctx, "*ESR?")
	if err != nil {
		return 0, err
	}

	return parseBit(res)
}

// SetServiceRequestEnable sets the value of the Service Request Enable register.
func (h *Handler) SetServiceRequestEnable(ctx context.Context, bit uint8) error {
	cmd := fmt.Sprintf("*SRE %d", bit)
	return h.ExecContext(ctx, cmd)
}

// QueryServiceRequestEnable queries the Service Request Enable.
func (h *Handler) QueryServiceRequestEnable(ctx context.Context) (bit uint8, err error) {
	res, err := h.QueryContext(ctx, "*SRE?")
	if err != nil {
		return 0, err
	}

	return parseBit(res)
}

// QueryStatusByteRegister queries the Status Byte Register.
func (h *Handler) QueryStatusByteRegister(ctx context.Context) (bit uint8, err error) {
	res, err := h.QueryContext(ctx, "*STB?")
	if err != nil {
		return 0, err
	}

	return parseBit(res)
}

// Recall restored the instrument to a state that was previously stored
// in locations 0 through 9 with the Save.
func (h *Handler) Recall(ctx context.Context, mem uint8) error {
	cmd := fmt.Sprintf("*RCL %d", mem)
	return h.ExecContext(ctx, cmd)
}

// Save saves the instrument setting to one of the ten non-volatile memory locations.
func (h *Handler) Save(ctx context.Context, mem uint8) error {
	cmd := fmt.Sprintf("*SAV %d", mem)
	return h.ExecContext(ctx, cmd)
}

var bitRegexp = regexp.MustCompile(`\+(\d+)`)

func parseBit(s string) (bit uint8, err error) {
	re := bitRegexp.Copy()
	g := re.FindStringSubmatch(s)
	if g == nil {
		return 0, fmt.Errorf("invalid bit format: %s", s)
	}
	n, err := strconv.ParseUint(g[1], 10, 8)
	if err != nil {
		return 0, err
	}
	bit = uint8(n)
	return bit, nil
}
