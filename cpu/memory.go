package cpu

// Bus astrae lo spazio memoria a 16 bit visto dal 6502.
type Bus interface {
	Read(addr uint16) byte
	Write(addr uint16, value byte)
}

// RAM e' una memoria piatta da 64 KB.
type RAM struct {
	Data [1 << 16]byte
}

// NewRAM crea una RAM azzerata.
func NewRAM() *RAM { return &RAM{} }

func (m *RAM) Read(addr uint16) byte { return m.Data[addr] }

func (m *RAM) Write(addr uint16, value byte) { m.Data[addr] = value }

// LoadAt copia data a partire da addr, con wrap naturale a 16 bit.
func (m *RAM) LoadAt(addr uint16, data []byte) {
	for i, b := range data {
		m.Write(addr+uint16(i), b)
	}
}
