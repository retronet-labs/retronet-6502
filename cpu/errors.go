package cpu

import "fmt"

// ErrNilBus segnala una CPU senza bus memoria collegato.
var ErrNilBus = fmt.Errorf("cpu6502: bus memoria nil")

// IllegalOpcodeError descrive un opcode non documentato/non implementato.
type IllegalOpcodeError struct {
	PC     uint16
	Opcode byte
}

func (e *IllegalOpcodeError) Error() string {
	return fmt.Sprintf("cpu6502: opcode illegale 0x%02X a 0x%04X", e.Opcode, e.PC)
}
