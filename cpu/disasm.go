package cpu

import "fmt"

// Disassemble restituisce il testo dell'istruzione a addr e la sua lunghezza.
func (c *CPU6502) Disassemble(addr uint16) (string, int) {
	if c.Mem == nil {
		return "???", 1
	}
	op := c.read(addr)
	def := opcodes[op]
	if def.op == opIllegal {
		return fmt.Sprintf(".byte $%02X", op), 1
	}
	lo := c.read(addr + 1)
	hi := c.read(addr + 2)
	word := uint16(lo) | uint16(hi)<<8
	switch def.mode {
	case modeImplied:
		return def.mnemonic, int(def.bytes)
	case modeAccumulator:
		return def.mnemonic + " A", int(def.bytes)
	case modeImmediate:
		return fmt.Sprintf("%s #$%02X", def.mnemonic, lo), int(def.bytes)
	case modeZeroPage:
		return fmt.Sprintf("%s $%02X", def.mnemonic, lo), int(def.bytes)
	case modeZeroPageX:
		return fmt.Sprintf("%s $%02X,X", def.mnemonic, lo), int(def.bytes)
	case modeZeroPageY:
		return fmt.Sprintf("%s $%02X,Y", def.mnemonic, lo), int(def.bytes)
	case modeAbsolute:
		return fmt.Sprintf("%s $%04X", def.mnemonic, word), int(def.bytes)
	case modeAbsoluteX:
		return fmt.Sprintf("%s $%04X,X", def.mnemonic, word), int(def.bytes)
	case modeAbsoluteY:
		return fmt.Sprintf("%s $%04X,Y", def.mnemonic, word), int(def.bytes)
	case modeIndirect:
		return fmt.Sprintf("%s ($%04X)", def.mnemonic, word), int(def.bytes)
	case modeIndexedIndirect:
		return fmt.Sprintf("%s ($%02X,X)", def.mnemonic, lo), int(def.bytes)
	case modeIndirectIndexed:
		return fmt.Sprintf("%s ($%02X),Y", def.mnemonic, lo), int(def.bytes)
	case modeRelative:
		target := uint16(int(addr+2) + int(int8(lo)))
		return fmt.Sprintf("%s $%04X", def.mnemonic, target), int(def.bytes)
	default:
		return def.mnemonic, int(def.bytes)
	}
}
