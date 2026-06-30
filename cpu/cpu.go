package cpu

import "github.com/retronet-labs/retronet-hardware/bridge/i6502"

const (
	vectorNMI   uint16 = 0xFFFA
	vectorReset uint16 = 0xFFFC
	vectorIRQ   uint16 = 0xFFFE
	stackBase   uint16 = 0x0100
)

// CPU6502 rappresenta lo stato visibile del MOS/NMOS 6502.
type CPU6502 struct {
	A  byte
	X  byte
	Y  byte
	SP byte
	PC uint16

	Carry            bool
	Zero             bool
	InterruptDisable bool
	Decimal          bool
	Overflow         bool
	Negative         bool

	Mem Bus

	InstructionCount uint64
	CycleCount       uint64
	LastCycles       int

	alu ALUBackend
}

// NewCPU6502 crea una CPU con RAM piatta da 64 KB e ALU a porte logiche.
func NewCPU6502() *CPU6502 {
	c := &CPU6502{Mem: NewRAM(), alu: Gate}
	c.Reset()
	return c
}

// NewCPU6502WithALU crea una CPU scegliendo il backend ALU.
func NewCPU6502WithALU(backend ALUBackend) *CPU6502 {
	c := NewCPU6502()
	c.SetALU(backend)
	return c
}

// SetALU sceglie il backend aritmetico-logico. nil ripristina Gate.
func (c *CPU6502) SetALU(backend ALUBackend) {
	if backend == nil {
		c.alu = Gate
		return
	}
	c.alu = backend
}

func (c *CPU6502) backend() ALUBackend {
	if c.alu == nil {
		return Gate
	}
	return c.alu
}

// Reset applica lo stato di reset del 6502: legge il vettore FFFC/FFFD, imposta
// SP a FD e abilita il mascheramento IRQ. La memoria non viene modificata.
func (c *CPU6502) Reset() {
	if c.Mem == nil {
		c.Mem = NewRAM()
	}
	c.A, c.X, c.Y = 0, 0, 0
	c.SP = 0xFD
	c.Carry, c.Zero = false, false
	c.InterruptDisable = true
	c.Decimal = false
	c.Overflow, c.Negative = false, false
	c.PC = c.read16(vectorReset)
	c.LastCycles = 7
}

// PackFlags restituisce il byte P. breakFlag controlla il bit B salvato da
// PHP/BRK; il bit 5 e' sempre a 1.
func (c *CPU6502) PackFlags(breakFlag bool) byte {
	var p byte = 0x20
	if c.Carry {
		p |= 0x01
	}
	if c.Zero {
		p |= 0x02
	}
	if c.InterruptDisable {
		p |= 0x04
	}
	if c.Decimal {
		p |= 0x08
	}
	if breakFlag {
		p |= 0x10
	}
	if c.Overflow {
		p |= 0x40
	}
	if c.Negative {
		p |= 0x80
	}
	return p
}

// SetFlagsByte ripristina i flag dal byte P. I bit 4 e 5 non sono stato reale.
func (c *CPU6502) SetFlagsByte(p byte) {
	c.Carry = p&0x01 != 0
	c.Zero = p&0x02 != 0
	c.InterruptDisable = p&0x04 != 0
	c.Decimal = p&0x08 != 0
	c.Overflow = p&0x40 != 0
	c.Negative = p&0x80 != 0
}

// IRQ consegna un interrupt mascherabile, se I=0.
func (c *CPU6502) IRQ() {
	if c.InterruptDisable {
		return
	}
	c.interrupt(vectorIRQ, false, 7)
}

// NMI consegna un interrupt non mascherabile.
func (c *CPU6502) NMI() { c.interrupt(vectorNMI, false, 7) }

func (c *CPU6502) read(addr uint16) byte { return c.Mem.Read(addr) }

func (c *CPU6502) write(addr uint16, value byte) { c.Mem.Write(addr, value) }

func (c *CPU6502) read16(addr uint16) uint16 {
	lo := uint16(c.read(addr))
	hi := uint16(c.read(addr + 1))
	return lo | hi<<8
}

func (c *CPU6502) read16Bug(addr uint16) uint16 {
	lo := uint16(c.read(addr))
	hiAddr := addr&0xFF00 | uint16(byte(addr+1))
	hi := uint16(c.read(hiAddr))
	return lo | hi<<8
}

func (c *CPU6502) fetch8() byte {
	v := c.read(c.PC)
	c.PC++
	return v
}

func (c *CPU6502) fetch16() uint16 {
	lo := uint16(c.fetch8())
	hi := uint16(c.fetch8())
	return lo | hi<<8
}

func (c *CPU6502) push(v byte) {
	c.write(stackBase|uint16(c.SP), v)
	c.SP--
}

func (c *CPU6502) pop() byte {
	c.SP++
	return c.read(stackBase | uint16(c.SP))
}

func (c *CPU6502) push16(v uint16) {
	c.push(byte(v >> 8))
	c.push(byte(v))
}

func (c *CPU6502) pop16() uint16 {
	lo := uint16(c.pop())
	hi := uint16(c.pop())
	return lo | hi<<8
}

func (c *CPU6502) setNZ(v byte) {
	c.Zero = v == 0
	c.Negative = v&0x80 != 0
}

func (c *CPU6502) applyNZ(f i6502.Flags) {
	c.Zero = f.Zero
	c.Negative = f.Negative
}

func (c *CPU6502) applyALUFlags(f i6502.Flags) {
	c.Carry = f.Carry
	c.Zero = f.Zero
	c.Negative = f.Negative
	c.Overflow = f.Overflow
}

func (c *CPU6502) interrupt(vector uint16, breakFlag bool, cycles int) {
	c.push16(c.PC)
	c.push(c.PackFlags(breakFlag))
	c.InterruptDisable = true
	c.PC = c.read16(vector)
	c.LastCycles = cycles
	c.CycleCount += uint64(cycles)
}

func samePage(a, b uint16) bool { return a&0xFF00 == b&0xFF00 }
