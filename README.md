# retronet-6502 - Emulatore MOS/NMOS 6502

Emulatore della CPU **MOS Technology 6502** scritto in Go, parte
dell'ecosistema RetroNet. Il core e' importabile, documentato in italiano e usa
due ALU intercambiabili:

- `cpu.Gate` (default): ALU a porte logiche tramite
  `retronet-hardware/bridge/i6502`;
- `cpu.Native`: stessa semantica con operatori Go, usata come oracolo veloce dei
  test differenziali.

Il target del primo rilascio e' il **6502 NMOS classico**: opcode documentati,
decimal mode, interrupt, stack page `$0100`, bug storico di `JMP ($xxFF)` e
conteggio cicli ufficiale aggregato.

## Quick Start

```bash
go test ./...
go run ./cmd/retronet-6502 -conformance
go run ./cmd/retronet-6502 -bin programma.bin -load 0x8000 -steps 1000
go run ./cmd/retronet-6502 -bin programma.bin -disasm 8
go run ./cmd/retronet-6502 -bin programma.bin -alu native -trace
```

Uso come libreria:

```go
c := cpu.NewCPU6502()        // RAM 64 KB, ALU Gate
ram := c.Mem.(*cpu.RAM)
ram.LoadAt(0x8000, []byte{0xA9, 0x2A, 0x8D, 0x00, 0x02}) // LDA #$2A; STA $0200
ram.Write(0xFFFC, 0x00)
ram.Write(0xFFFD, 0x80)
c.Reset()
c.Run(2)
```

## Stato

- Registri `A`, `X`, `Y`, `SP`, `PC` e flag `N V - B D I Z C`.
- Memoria piatta da 64 KB (`cpu.RAM`) dietro interfaccia `cpu.Bus`.
- Reset vector `$FFFC/$FFFD`, NMI `$FFFA/$FFFB`, IRQ/BRK `$FFFE/$FFFF`.
- Stack hardware su pagina `$0100`, con `PHP/PLP`, `PHA/PLA`, `JSR/RTS`,
  `BRK/RTI`.
- Opcode documentati NMOS completi; opcode non documentati restituiscono
  `IllegalOpcodeError`.
- Addressing mode ufficiali, inclusi indexed-indirect, indirect-indexed,
  relative branch e bug `JMP ($xxFF)`.
- Decimal mode per `ADC`/`SBC`; per il modello NMOS `Z/N/V` derivano dal
  risultato binario pre-correzione BCD.
- Disassembler minimale usato da `-trace` e `-disasm`.
- Conformance sintetica interna, senza ROM storiche vendorizzate.

## Documentazione

- [Architettura](docs/architettura.md)
- [Istruzioni e addressing mode](docs/istruzioni.md)
- [Flag e decimal mode](docs/flags.md)
- [Timing](docs/timing.md)
- [Conformance](docs/conformance.md)

## Sviluppo Locale

Un clone pulito compila dalle versioni taggate di `retronet-hardware` e
`retronet-logic`. Per co-sviluppare i tre moduli come sibling:

```bash
go work init . ../retronet-hardware ../retronet-logic
go test ./...
go vet ./...
```

`go.work` e' locale e non va versionato.

## Limiti

- Niente 65C02.
- Niente opcode illegali NMOS.
- Timing a cicli aggregati, non segnali pin-by-pin o bus-cycle.
- Nessuna ROM storica inclusa.

## Licenza

MIT.
