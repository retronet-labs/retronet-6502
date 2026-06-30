package cpu

import "github.com/retronet-labs/retronet-hardware/bridge/i6502"

type nativeBackend struct{}

func (nativeBackend) ADC(a, value byte, carryIn bool, decimal bool) (byte, i6502.Flags) {
	c := 0
	if carryIn {
		c = 1
	}
	sum := int(a) + int(value) + c
	bin := byte(sum)
	flags := i6502.Flags{
		Carry:    sum > 0xFF,
		Zero:     bin == 0,
		Negative: bin&0x80 != 0,
		Overflow: (^(a ^ value) & (a ^ bin) & 0x80) != 0,
	}
	if !decimal {
		return bin, flags
	}
	result := sum
	if int(a&0x0F)+int(value&0x0F)+c > 9 {
		result += 0x06
	}
	if sum > 0x99 {
		result += 0x60
		flags.Carry = true
	} else {
		flags.Carry = false
	}
	return byte(result), flags
}

func (nativeBackend) SBC(a, value byte, carryIn bool, decimal bool) (byte, i6502.Flags) {
	borrow := 1
	if carryIn {
		borrow = 0
	}
	diff := int(a) - int(value) - borrow
	bin := byte(diff)
	flags := i6502.Flags{
		Carry:    diff >= 0,
		Zero:     bin == 0,
		Negative: bin&0x80 != 0,
		Overflow: ((a ^ value) & (a ^ bin) & 0x80) != 0,
	}
	if !decimal {
		return bin, flags
	}
	result := diff
	if int(a&0x0F)-borrow < int(value&0x0F) {
		result -= 0x06
	}
	if diff < 0 {
		result -= 0x60
	}
	return byte(result), flags
}

func (nativeBackend) Compare(reg, value byte) (byte, i6502.Flags) {
	out := byte(reg - value)
	return out, i6502.Flags{Carry: reg >= value, Zero: out == 0, Negative: out&0x80 != 0}
}

func (nativeBackend) Logic(op i6502.Op, a, value byte) (byte, i6502.Flags) {
	var out byte
	switch op {
	case i6502.OpAND:
		out = a & value
	case i6502.OpEOR:
		out = a ^ value
	default:
		out = a | value
	}
	return out, i6502.Flags{Zero: out == 0, Negative: out&0x80 != 0}
}

func (n nativeBackend) BIT(a, value byte) i6502.Flags {
	_, f := n.Logic(i6502.OpAND, a, value)
	return i6502.Flags{Zero: f.Zero, Negative: value&0x80 != 0, Overflow: value&0x40 != 0}
}

func (nativeBackend) Increment(value byte) (byte, i6502.Flags) {
	out := value + 1
	return out, i6502.Flags{Zero: out == 0, Negative: out&0x80 != 0}
}

func (nativeBackend) Decrement(value byte) (byte, i6502.Flags) {
	out := value - 1
	return out, i6502.Flags{Zero: out == 0, Negative: out&0x80 != 0}
}

func (nativeBackend) ShiftLeft(value byte) (byte, i6502.Flags) {
	out := value << 1
	return out, i6502.Flags{Carry: value&0x80 != 0, Zero: out == 0, Negative: out&0x80 != 0}
}

func (nativeBackend) ShiftRight(value byte) (byte, i6502.Flags) {
	out := value >> 1
	return out, i6502.Flags{Carry: value&0x01 != 0, Zero: out == 0, Negative: false}
}

func (nativeBackend) RotateLeft(value byte, carryIn bool) (byte, i6502.Flags) {
	out := value << 1
	if carryIn {
		out |= 1
	}
	return out, i6502.Flags{Carry: value&0x80 != 0, Zero: out == 0, Negative: out&0x80 != 0}
}

func (nativeBackend) RotateRight(value byte, carryIn bool) (byte, i6502.Flags) {
	out := value >> 1
	if carryIn {
		out |= 0x80
	}
	return out, i6502.Flags{Carry: value&0x01 != 0, Zero: out == 0, Negative: out&0x80 != 0}
}
