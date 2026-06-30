package cpu

import "github.com/retronet-labs/retronet-hardware/bridge/i6502"

// ALUBackend astrae il motore aritmetico-logico del 6502. Gate e' il default:
// inoltra al bridge i6502 costruito sopra retronet-logic. Native usa operatori
// Go con la stessa semantica ed e' l'oracolo veloce dei test differenziali.
type ALUBackend interface {
	ADC(a, value byte, carryIn bool, decimal bool) (byte, i6502.Flags)
	SBC(a, value byte, carryIn bool, decimal bool) (byte, i6502.Flags)
	Compare(reg, value byte) (byte, i6502.Flags)
	Logic(op i6502.Op, a, value byte) (byte, i6502.Flags)
	BIT(a, value byte) i6502.Flags
	Increment(value byte) (byte, i6502.Flags)
	Decrement(value byte) (byte, i6502.Flags)
	ShiftLeft(value byte) (byte, i6502.Flags)
	ShiftRight(value byte) (byte, i6502.Flags)
	RotateLeft(value byte, carryIn bool) (byte, i6502.Flags)
	RotateRight(value byte, carryIn bool) (byte, i6502.Flags)
}

// Gate e' il backend ALU a porte logiche (default).
var Gate ALUBackend = gateBackend{}

// Native e' il backend ALU con operatori Go.
var Native ALUBackend = nativeBackend{}

type gateBackend struct{}

func (gateBackend) ADC(a, value byte, carryIn bool, decimal bool) (byte, i6502.Flags) {
	return i6502.ADC(a, value, carryIn, decimal)
}

func (gateBackend) SBC(a, value byte, carryIn bool, decimal bool) (byte, i6502.Flags) {
	return i6502.SBC(a, value, carryIn, decimal)
}

func (gateBackend) Compare(reg, value byte) (byte, i6502.Flags) { return i6502.Compare(reg, value) }

func (gateBackend) Logic(op i6502.Op, a, value byte) (byte, i6502.Flags) {
	return i6502.Logic(op, a, value)
}

func (gateBackend) BIT(a, value byte) i6502.Flags { return i6502.BIT(a, value) }

func (gateBackend) Increment(value byte) (byte, i6502.Flags) { return i6502.Increment(value) }

func (gateBackend) Decrement(value byte) (byte, i6502.Flags) { return i6502.Decrement(value) }

func (gateBackend) ShiftLeft(value byte) (byte, i6502.Flags) { return i6502.ShiftLeft(value) }

func (gateBackend) ShiftRight(value byte) (byte, i6502.Flags) { return i6502.ShiftRight(value) }

func (gateBackend) RotateLeft(value byte, carryIn bool) (byte, i6502.Flags) {
	return i6502.RotateLeft(value, carryIn)
}

func (gateBackend) RotateRight(value byte, carryIn bool) (byte, i6502.Flags) {
	return i6502.RotateRight(value, carryIn)
}
