package cpu

import (
	"testing"
)

func newTestCPU(code []byte) *CPU6502 {
	c := NewCPU6502()
	ram := c.Mem.(*RAM)
	ram.LoadAt(0x8000, code)
	ram.Write(vectorReset, 0x00)
	ram.Write(vectorReset+1, 0x80)
	ram.Write(vectorIRQ, 0x00)
	ram.Write(vectorIRQ+1, 0x90)
	c.Reset()
	return c
}

func stepN(t *testing.T, c *CPU6502, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		if err := c.Step(); err != nil {
			t.Fatalf("step %d: %v", i, err)
		}
	}
}

func TestResetReadsVectorAndInitializesState(t *testing.T) {
	c := NewCPU6502()
	ram := c.Mem.(*RAM)
	ram.Write(vectorReset, 0x34)
	ram.Write(vectorReset+1, 0x12)
	c.A, c.X, c.Y = 1, 2, 3
	c.Reset()
	if c.PC != 0x1234 || c.SP != 0xFD || !c.InterruptDisable || c.Decimal {
		t.Fatalf("stato reset errato: PC=%04X SP=%02X I=%v D=%v", c.PC, c.SP, c.InterruptDisable, c.Decimal)
	}
	if c.A != 0 || c.X != 0 || c.Y != 0 {
		t.Fatalf("registri non azzerati: A=%02X X=%02X Y=%02X", c.A, c.X, c.Y)
	}
}

func TestLoadStoreAndZeroPageWrap(t *testing.T) {
	c := newTestCPU([]byte{
		0xA2, 0x01, // LDX #1
		0xA9, 0x42, // LDA #$42
		0x95, 0xFF, // STA $FF,X -> $00
		0xB5, 0xFF, // LDA $FF,X -> $00
	})
	stepN(t, c, 4)
	if got := c.Mem.Read(0x0000); got != 0x42 {
		t.Fatalf("mem[0000]=%02X, atteso 42", got)
	}
	if c.A != 0x42 || c.Zero || c.Negative {
		t.Fatalf("A/flag errati: A=%02X Z=%v N=%v", c.A, c.Zero, c.Negative)
	}
}

func TestDecimalADCAndSBC(t *testing.T) {
	c := newTestCPU([]byte{
		0xF8,       // SED
		0x18,       // CLC
		0xA9, 0x45, // LDA #$45
		0x69, 0x55, // ADC #$55 -> $00 C=1 (BCD 100)
		0xE9, 0x01, // SBC #$01 con C=1 -> $99
	})
	stepN(t, c, 5)
	if c.A != 0x99 || c.Carry || c.Zero {
		t.Fatalf("decimal result A=%02X C=%v Z=%v", c.A, c.Carry, c.Zero)
	}
}

func TestJSRRTSAndStack(t *testing.T) {
	c := newTestCPU([]byte{
		0x20, 0x06, 0x80, // JSR sub
		0x8D, 0x00, 0x02, // STA $0200
		0xA9, 0x77, // sub: LDA #$77
		0x60, // RTS
	})
	stepN(t, c, 4)
	if c.Mem.Read(0x0200) != 0x77 || c.SP != 0xFD {
		t.Fatalf("JSR/RTS errato: mem=%02X SP=%02X PC=%04X", c.Mem.Read(0x0200), c.SP, c.PC)
	}
}

func TestBranchCyclesAndPageCrossing(t *testing.T) {
	c := newTestCPU([]byte{0xD0, 0x02}) // BNE +2
	c.PC = 0x80FD
	c.Zero = false
	c.Mem.Write(0x80FD, 0xD0)
	c.Mem.Write(0x80FE, 0x02)
	if err := c.Step(); err != nil {
		t.Fatal(err)
	}
	if c.PC != 0x8101 || c.LastCycles != 4 {
		t.Fatalf("branch PC=%04X cycles=%d", c.PC, c.LastCycles)
	}
}

func TestJMPIndirectPageBug(t *testing.T) {
	c := newTestCPU([]byte{0x6C, 0xFF, 0x12}) // JMP ($12FF)
	c.Mem.Write(0x12FF, 0x34)
	c.Mem.Write(0x1200, 0x56)
	c.Mem.Write(0x1300, 0x99)
	if err := c.Step(); err != nil {
		t.Fatal(err)
	}
	if c.PC != 0x5634 {
		t.Fatalf("JMP indirect bug PC=%04X, atteso 5634", c.PC)
	}
}

func TestBRKPushesBreakFlagAndJumpsIRQVector(t *testing.T) {
	c := newTestCPU([]byte{0x00, 0xEA})
	if err := c.Step(); err != nil {
		t.Fatal(err)
	}
	if c.PC != 0x9000 || c.SP != 0xFA || !c.InterruptDisable {
		t.Fatalf("BRK state PC=%04X SP=%02X I=%v", c.PC, c.SP, c.InterruptDisable)
	}
	p := c.Mem.Read(0x01FB)
	if p&0x10 == 0 || p&0x20 == 0 {
		t.Fatalf("flags pushed=%02X, attesi B e bit5", p)
	}
}

func TestOfficialOpcodeCount(t *testing.T) {
	n := 0
	for _, d := range opcodes {
		if d.op != opIllegal {
			n++
		}
	}
	if n != 151 {
		t.Fatalf("opcode documentati = %d, attesi 151", n)
	}
}

func TestGateVsNativeALUDifferential(t *testing.T) {
	for a := 0; a <= 0xFF; a++ {
		for v := 0; v <= 0xFF; v++ {
			for _, carry := range []bool{false, true} {
				for _, dec := range []bool{false, true} {
					gr, gf := Gate.ADC(byte(a), byte(v), carry, dec)
					nr, nf := Native.ADC(byte(a), byte(v), carry, dec)
					if gr != nr || gf != nf {
						t.Fatalf("ADC a=%02X v=%02X c=%v d=%v gate=%02X %+v native=%02X %+v", a, v, carry, dec, gr, gf, nr, nf)
					}
					gr, gf = Gate.SBC(byte(a), byte(v), carry, dec)
					nr, nf = Native.SBC(byte(a), byte(v), carry, dec)
					if gr != nr || gf != nf {
						t.Fatalf("SBC a=%02X v=%02X c=%v d=%v gate=%02X %+v native=%02X %+v", a, v, carry, dec, gr, gf, nr, nf)
					}
				}
			}
		}
	}
}
