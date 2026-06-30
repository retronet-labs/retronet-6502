package cpu

type addrMode byte

const (
	modeImplied addrMode = iota
	modeAccumulator
	modeImmediate
	modeZeroPage
	modeZeroPageX
	modeZeroPageY
	modeAbsolute
	modeAbsoluteX
	modeAbsoluteY
	modeIndirect
	modeIndexedIndirect
	modeIndirectIndexed
	modeRelative
)

type operation byte

const (
	opIllegal operation = iota
	opADC
	opAND
	opASL
	opBCC
	opBCS
	opBEQ
	opBIT
	opBMI
	opBNE
	opBPL
	opBRK
	opBVC
	opBVS
	opCLC
	opCLD
	opCLI
	opCLV
	opCMP
	opCPX
	opCPY
	opDEC
	opDEX
	opDEY
	opEOR
	opINC
	opINX
	opINY
	opJMP
	opJSR
	opLDA
	opLDX
	opLDY
	opLSR
	opNOP
	opORA
	opPHA
	opPHP
	opPLA
	opPLP
	opROL
	opROR
	opRTI
	opRTS
	opSBC
	opSEC
	opSED
	opSEI
	opSTA
	opSTX
	opSTY
	opTAX
	opTAY
	opTSX
	opTXA
	opTXS
	opTYA
)

type opcodeDef struct {
	mnemonic  string
	op        operation
	mode      addrMode
	bytes     byte
	cycles    byte
	pageExtra bool
}

var opcodes [256]opcodeDef

func init() {
	for i := range opcodes {
		opcodes[i] = opcodeDef{mnemonic: "???", op: opIllegal, bytes: 1}
	}
	set(0x00, "BRK", opBRK, modeImplied, 1, 7, false)
	set(0xEA, "NOP", opNOP, modeImplied, 1, 2, false)

	// ORA
	set(0x09, "ORA", opORA, modeImmediate, 2, 2, false)
	set(0x05, "ORA", opORA, modeZeroPage, 2, 3, false)
	set(0x15, "ORA", opORA, modeZeroPageX, 2, 4, false)
	set(0x0D, "ORA", opORA, modeAbsolute, 3, 4, false)
	set(0x1D, "ORA", opORA, modeAbsoluteX, 3, 4, true)
	set(0x19, "ORA", opORA, modeAbsoluteY, 3, 4, true)
	set(0x01, "ORA", opORA, modeIndexedIndirect, 2, 6, false)
	set(0x11, "ORA", opORA, modeIndirectIndexed, 2, 5, true)

	// AND/EOR/ADC/SBC/CMP condividono gli addressing mode di ORA.
	addALU("AND", opAND, 0x29, 0x25, 0x35, 0x2D, 0x3D, 0x39, 0x21, 0x31)
	addALU("EOR", opEOR, 0x49, 0x45, 0x55, 0x4D, 0x5D, 0x59, 0x41, 0x51)
	addALU("ADC", opADC, 0x69, 0x65, 0x75, 0x6D, 0x7D, 0x79, 0x61, 0x71)
	addALU("SBC", opSBC, 0xE9, 0xE5, 0xF5, 0xED, 0xFD, 0xF9, 0xE1, 0xF1)
	addALU("CMP", opCMP, 0xC9, 0xC5, 0xD5, 0xCD, 0xDD, 0xD9, 0xC1, 0xD1)

	// ASL/ROL/LSR/ROR
	addRMW("ASL", opASL, 0x0A, 0x06, 0x16, 0x0E, 0x1E)
	addRMW("ROL", opROL, 0x2A, 0x26, 0x36, 0x2E, 0x3E)
	addRMW("LSR", opLSR, 0x4A, 0x46, 0x56, 0x4E, 0x5E)
	addRMW("ROR", opROR, 0x6A, 0x66, 0x76, 0x6E, 0x7E)

	// Store
	set(0x85, "STA", opSTA, modeZeroPage, 2, 3, false)
	set(0x95, "STA", opSTA, modeZeroPageX, 2, 4, false)
	set(0x8D, "STA", opSTA, modeAbsolute, 3, 4, false)
	set(0x9D, "STA", opSTA, modeAbsoluteX, 3, 5, false)
	set(0x99, "STA", opSTA, modeAbsoluteY, 3, 5, false)
	set(0x81, "STA", opSTA, modeIndexedIndirect, 2, 6, false)
	set(0x91, "STA", opSTA, modeIndirectIndexed, 2, 6, false)
	set(0x84, "STY", opSTY, modeZeroPage, 2, 3, false)
	set(0x94, "STY", opSTY, modeZeroPageX, 2, 4, false)
	set(0x8C, "STY", opSTY, modeAbsolute, 3, 4, false)
	set(0x86, "STX", opSTX, modeZeroPage, 2, 3, false)
	set(0x96, "STX", opSTX, modeZeroPageY, 2, 4, false)
	set(0x8E, "STX", opSTX, modeAbsolute, 3, 4, false)

	// Load
	addLoad("LDA", opLDA, 0xA9, 0xA5, 0xB5, 0xAD, 0xBD, 0xB9, 0xA1, 0xB1)
	set(0xA2, "LDX", opLDX, modeImmediate, 2, 2, false)
	set(0xA6, "LDX", opLDX, modeZeroPage, 2, 3, false)
	set(0xB6, "LDX", opLDX, modeZeroPageY, 2, 4, false)
	set(0xAE, "LDX", opLDX, modeAbsolute, 3, 4, false)
	set(0xBE, "LDX", opLDX, modeAbsoluteY, 3, 4, true)
	set(0xA0, "LDY", opLDY, modeImmediate, 2, 2, false)
	set(0xA4, "LDY", opLDY, modeZeroPage, 2, 3, false)
	set(0xB4, "LDY", opLDY, modeZeroPageX, 2, 4, false)
	set(0xAC, "LDY", opLDY, modeAbsolute, 3, 4, false)
	set(0xBC, "LDY", opLDY, modeAbsoluteX, 3, 4, true)

	// Compare speciali.
	set(0xE0, "CPX", opCPX, modeImmediate, 2, 2, false)
	set(0xE4, "CPX", opCPX, modeZeroPage, 2, 3, false)
	set(0xEC, "CPX", opCPX, modeAbsolute, 3, 4, false)
	set(0xC0, "CPY", opCPY, modeImmediate, 2, 2, false)
	set(0xC4, "CPY", opCPY, modeZeroPage, 2, 3, false)
	set(0xCC, "CPY", opCPY, modeAbsolute, 3, 4, false)
	set(0x24, "BIT", opBIT, modeZeroPage, 2, 3, false)
	set(0x2C, "BIT", opBIT, modeAbsolute, 3, 4, false)

	// INC/DEC memoria.
	set(0xE6, "INC", opINC, modeZeroPage, 2, 5, false)
	set(0xF6, "INC", opINC, modeZeroPageX, 2, 6, false)
	set(0xEE, "INC", opINC, modeAbsolute, 3, 6, false)
	set(0xFE, "INC", opINC, modeAbsoluteX, 3, 7, false)
	set(0xC6, "DEC", opDEC, modeZeroPage, 2, 5, false)
	set(0xD6, "DEC", opDEC, modeZeroPageX, 2, 6, false)
	set(0xCE, "DEC", opDEC, modeAbsolute, 3, 6, false)
	set(0xDE, "DEC", opDEC, modeAbsoluteX, 3, 7, false)

	// Branch.
	set(0x10, "BPL", opBPL, modeRelative, 2, 2, false)
	set(0x30, "BMI", opBMI, modeRelative, 2, 2, false)
	set(0x50, "BVC", opBVC, modeRelative, 2, 2, false)
	set(0x70, "BVS", opBVS, modeRelative, 2, 2, false)
	set(0x90, "BCC", opBCC, modeRelative, 2, 2, false)
	set(0xB0, "BCS", opBCS, modeRelative, 2, 2, false)
	set(0xD0, "BNE", opBNE, modeRelative, 2, 2, false)
	set(0xF0, "BEQ", opBEQ, modeRelative, 2, 2, false)

	// Salti, stack, flag, trasferimenti.
	set(0x4C, "JMP", opJMP, modeAbsolute, 3, 3, false)
	set(0x6C, "JMP", opJMP, modeIndirect, 3, 5, false)
	set(0x20, "JSR", opJSR, modeAbsolute, 3, 6, false)
	set(0x60, "RTS", opRTS, modeImplied, 1, 6, false)
	set(0x40, "RTI", opRTI, modeImplied, 1, 6, false)
	set(0x48, "PHA", opPHA, modeImplied, 1, 3, false)
	set(0x08, "PHP", opPHP, modeImplied, 1, 3, false)
	set(0x68, "PLA", opPLA, modeImplied, 1, 4, false)
	set(0x28, "PLP", opPLP, modeImplied, 1, 4, false)
	set(0x18, "CLC", opCLC, modeImplied, 1, 2, false)
	set(0x38, "SEC", opSEC, modeImplied, 1, 2, false)
	set(0x58, "CLI", opCLI, modeImplied, 1, 2, false)
	set(0x78, "SEI", opSEI, modeImplied, 1, 2, false)
	set(0xB8, "CLV", opCLV, modeImplied, 1, 2, false)
	set(0xD8, "CLD", opCLD, modeImplied, 1, 2, false)
	set(0xF8, "SED", opSED, modeImplied, 1, 2, false)
	set(0xAA, "TAX", opTAX, modeImplied, 1, 2, false)
	set(0xA8, "TAY", opTAY, modeImplied, 1, 2, false)
	set(0xBA, "TSX", opTSX, modeImplied, 1, 2, false)
	set(0x8A, "TXA", opTXA, modeImplied, 1, 2, false)
	set(0x9A, "TXS", opTXS, modeImplied, 1, 2, false)
	set(0x98, "TYA", opTYA, modeImplied, 1, 2, false)
	set(0xE8, "INX", opINX, modeImplied, 1, 2, false)
	set(0xC8, "INY", opINY, modeImplied, 1, 2, false)
	set(0xCA, "DEX", opDEX, modeImplied, 1, 2, false)
	set(0x88, "DEY", opDEY, modeImplied, 1, 2, false)
}

func set(code byte, mnemonic string, op operation, mode addrMode, bytes, cycles byte, pageExtra bool) {
	opcodes[code] = opcodeDef{mnemonic: mnemonic, op: op, mode: mode, bytes: bytes, cycles: cycles, pageExtra: pageExtra}
}

func addALU(m string, op operation, imm, zp, zpx, abs, absx, absy, indx, indy byte) {
	set(imm, m, op, modeImmediate, 2, 2, false)
	set(zp, m, op, modeZeroPage, 2, 3, false)
	set(zpx, m, op, modeZeroPageX, 2, 4, false)
	set(abs, m, op, modeAbsolute, 3, 4, false)
	set(absx, m, op, modeAbsoluteX, 3, 4, true)
	set(absy, m, op, modeAbsoluteY, 3, 4, true)
	set(indx, m, op, modeIndexedIndirect, 2, 6, false)
	set(indy, m, op, modeIndirectIndexed, 2, 5, true)
}

func addLoad(m string, op operation, imm, zp, zpx, abs, absx, absy, indx, indy byte) {
	addALU(m, op, imm, zp, zpx, abs, absx, absy, indx, indy)
}

func addRMW(m string, op operation, acc, zp, zpx, abs, absx byte) {
	set(acc, m, op, modeAccumulator, 1, 2, false)
	set(zp, m, op, modeZeroPage, 2, 5, false)
	set(zpx, m, op, modeZeroPageX, 2, 6, false)
	set(abs, m, op, modeAbsolute, 3, 6, false)
	set(absx, m, op, modeAbsoluteX, 3, 7, false)
}
