package scpi

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/xerrors"
)

// Handler is a handler for a device controlled using SCPI commands.
type Handler struct {
	Client
}

// Reset resets the instrument to a factory pre-defined condition and clears the error log.
func (h *Handler) Reset() error {
	return h.Exec("*RST;*CLS")
}

// WaitForComplete waits for all queued operations to complete up to the specified timeout.
func (h *Handler) WaitForComplete(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ch := make(chan error, 1)
	go func() {
		_, err := h.QueryContext(ctx, "*WAI;*OPC?")
		ch <- err
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		// TODO: Refactor the timeout error
		return xerrors.New("timeout")
	}
}

// Trigger triggers the device if, and only if,
// Bus Triggering is the type of trigger event selected.
// Otherwise, this command is ignored.
func (h *Handler) Trigger() error {
	return h.Exec("*TRG")
}

// Identify returns the identification data.
// The standards order is follows:
// - Manufacturer
// - Model number
// - Serial number (or 0)
// - Firmware version
func (h *Handler) Identify() (id string, err error) {
	res, err := h.Query("*IDN?")
	if err != nil {
		return "", nil
	}

	id = string(res)
	return id, nil
}

// SetEventStatusEnable sets the value in the enable register for the Standard Event Status group.
// The selected bits are then reported to bit 5 of the Status Byte.
func (h *Handler) SetEventStatusEnable(bits uint8) error {
	cmd := fmt.Sprintf("*ESE %d", bits)
	return h.Exec(cmd)
}

// QueryEventStatusEnable queries the event status enable.
func (h *Handler) QueryEventStatusEnable() (bits uint8, err error) {
	res, err := h.Query("*ESE?")
	if err != nil {
		return 0, err
	}

	return h.bytesToUint8(res)
}

// QueryEventStatusRegister queries the event status register.
// The register is cleared when it is executed.
func (h *Handler) QueryEventStatusRegister() (bits uint8, err error) {
	res, err := h.Query("*ESR?")
	if err != nil {
		return 0, err
	}

	return h.bytesToUint8(res)
}

// SetServiceRequestEnable sets the value of the Service Request Enable register.
func (h *Handler) SetServiceRequestEnable(bits uint8) error {
	cmd := fmt.Sprintf("*SRE %d", bits)
	return h.Exec(cmd)
}

// QueryServiceRequestEnable queries the Service Request Enable.
func (h *Handler) QueryServiceRequestEnable() (bits uint8, err error) {
	res, err := h.Query("*SRE?")
	if err != nil {
		return 0, err
	}

	return h.bytesToUint8(res)
}

// QueryStatusByteRegister queries the Status Byte Register.
func (h *Handler) QueryStatusByteRegister() (bits uint8, err error) {
	res, err := h.Query("*STB?")
	if err != nil {
		return 0, err
	}

	return h.bytesToUint8(res)
}

// Recall restored the instrument to a state that was previously stored
// in locations 0 through 9 with the Save.
func (h *Handler) Recall(mem uint8) error {
	if mem > 9 {
		// TODO: Refactor the timeout error
		return xerrors.New("only 0 to 10 are allowed")
	}

	cmd := fmt.Sprintf("*RCL %d", mem)
	return h.Exec(cmd)
}

// Save saves the instrument setting to one of the ten non-volatile memory locations.
func (h *Handler) Save(mem uint8) error {
	if mem > 9 {
		// TODO: Refactor the timeout error
		return xerrors.New("only 0 to 10 are allowed")
	}

	cmd := fmt.Sprintf("*SAV %d", mem)
	return h.Exec(cmd)
}

func (h *Handler) bytesToUint8(bytes []byte) (n uint8, err error) {
	num64, err := strconv.ParseUint(string(bytes), 10, 8)
	if err != nil {
		return 0, err
	}

	n = uint8(num64)
	return n, nil
}
