package scpi

import (
	"time"
)

// Handler is a handler for a device controlled using SCPI commands.
type Handler interface {
	// Reset resets the instrument to a factory pre-defined condition and clears the error log.
	Reset() error

	// WaitForComplete waits for all queued operations to complete up to the specified timeout.
	WaitForComplete(timeout time.Duration) error

	// Trigger triggers the device if, and only if,
	// Bus Triggering is the type of trigger event selected.
	// Otherwise, this command is ignored.
	Trigger() error

	// Identify returns the identification data.
	// The standards order is follows:
	// - Manufacturer
	// - Model number
	// - Serial number (or 0)
	// - Firmware version
	Identify() (id string, err error)

	// SetEventStatusEnable sets the value in the enable register for the Standard Event Status group.
	// The selected bits are then reported to bit 5 of the Status Byte.
	SetEventStatusEnable(bits uint8) error

	// QueryEventStatusEnable queries the event status enable.
	QueryEventStatusEnable() (bits uint8, err error)

	// QueryEventStatusRegister queries the event status register.
	// The register is cleared when it is executed.
	QueryEventStatusRegister() (bits uint8, err error)

	// SetServiceRequestEnable sets the value of the Service Request Enable register.
	SetServiceRequestEnable() (bits uint8, err error)

	// QueryServiceRequestEnable queries the Service Request Enable.
	QueryServiceRequestEnable() (bits uint8, err error)

	// QueryStatusByteRegister queries the Status Byte Register.
	QueryStatusByteRegister() (bits uint8, err error)

	// Recall restored the instrument to a state that was previously stored
	// in locations 0 through 9 with the Save.
	Recall(mem uint8) error

	// Save saves the instrument setting to one of the ten non-volatile memory locations.
	Save(mem uint8) error
}
