// Package conformance esegue piccoli programmi auto-verificanti sul core 6502,
// senza ROM storiche o dataset esterni.
package conformance

import (
	"fmt"

	"github.com/retronet-labs/retronet-6502/cpu"
)

// Case e' l'esito di un programma di prova.
type Case struct {
	Name   string
	OK     bool
	Detail string
}

// Result raccoglie gli esiti.
type Result struct {
	Cases []Case
}

func (r Result) Passed() int {
	n := 0
	for _, c := range r.Cases {
		if c.OK {
			n++
		}
	}
	return n
}

func (r Result) Failed() int { return len(r.Cases) - r.Passed() }

type program struct {
	name   string
	code   []byte
	steps  int
	verify func(*cpu.CPU6502) (bool, string)
}

var programs = []program{
	{
		name: "loop: 5 * 3 = 15",
		code: []byte{
			0xA2, 0x05, // LDX #5
			0xA9, 0x00, // LDA #0
			0x18,       // loop: CLC
			0x69, 0x03, // ADC #3
			0xCA,       // DEX
			0xD0, 0xFA, // BNE loop
			0x8D, 0x00, 0x02, // STA $0200
		},
		steps: 23,
		verify: func(c *cpu.CPU6502) (bool, string) {
			return c.Mem.Read(0x0200) == 15, fmt.Sprintf("$0200=%02X, atteso 0F", c.Mem.Read(0x0200))
		},
	},
	{
		name: "subroutine JSR/RTS",
		code: []byte{
			0x20, 0x07, 0x80, // JSR sub
			0x8D, 0x01, 0x02, // STA $0201
			0xEA,       // NOP
			0xA9, 0x42, // sub: LDA #$42
			0x60, // RTS
		},
		steps: 4,
		verify: func(c *cpu.CPU6502) (bool, string) {
			return c.Mem.Read(0x0201) == 0x42 && c.SP == 0xFD,
				fmt.Sprintf("$0201=%02X SP=%02X", c.Mem.Read(0x0201), c.SP)
		},
	},
	{
		name: "BCD: 45 + 55 = 100",
		code: []byte{
			0xF8,       // SED
			0x18,       // CLC
			0xA9, 0x45, // LDA #$45
			0x69, 0x55, // ADC #$55
			0x8D, 0x02, 0x02, // STA $0202
		},
		steps: 5,
		verify: func(c *cpu.CPU6502) (bool, string) {
			return c.Mem.Read(0x0202) == 0x00 && c.Carry,
				fmt.Sprintf("$0202=%02X C=%v", c.Mem.Read(0x0202), c.Carry)
		},
	},
	{
		name:  "JMP indirect page-wrap bug",
		code:  []byte{0x6C, 0xFF, 0x12},
		steps: 1,
		verify: func(c *cpu.CPU6502) (bool, string) {
			return c.PC == 0x5634, fmt.Sprintf("PC=%04X, atteso 5634", c.PC)
		},
	},
}

// Run esegue tutti i programmi col backend dato (Gate se nil).
func Run(backend cpu.ALUBackend) Result {
	if backend == nil {
		backend = cpu.Gate
	}
	var out Result
	for _, p := range programs {
		c := cpu.NewCPU6502WithALU(backend)
		ram := c.Mem.(*cpu.RAM)
		ram.LoadAt(0x8000, p.code)
		ram.Write(0xFFFC, 0x00)
		ram.Write(0xFFFD, 0x80)
		ram.Write(0x12FF, 0x34)
		ram.Write(0x1200, 0x56)
		c.Reset()
		if _, err := c.Run(p.steps); err != nil {
			out.Cases = append(out.Cases, Case{Name: p.name, OK: false, Detail: err.Error()})
			continue
		}
		ok, detail := p.verify(c)
		out.Cases = append(out.Cases, Case{Name: p.name, OK: ok, Detail: detail})
	}
	return out
}
