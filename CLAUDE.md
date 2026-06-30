# CLAUDE.md - retronet-6502

Emulatore MOS/NMOS 6502 in Go, coerente con gli altri core RetroNet: core
importabile, CLI, ALU `Gate`/`Native`, test e documentazione in italiano.

## Setup

```sh
go work init . ../retronet-hardware ../retronet-logic
go test ./...
go vet ./...
go run ./cmd/retronet-6502 -conformance
```

`go.work` e' locale e non versionato. Un clone pulito deve compilare dalle
versioni taggate di `retronet-hardware` e `retronet-logic`.

## Decisioni Da Preservare

- Target: NMOS 6502, non 65C02.
- Opcode illegali fuori scope: devono restituire `IllegalOpcodeError`.
- `cpu.Gate` e' il default e usa `retronet-hardware/bridge/i6502`.
- Decimal mode NMOS: `A` e `C` corretti BCD, `Z/N/V` dal risultato binario.
- Timing aggregato ufficiale, non bus-cycle.
- Nessuna ROM storica vendorizzata.

## Componenti

- `cpu/`: stato, memoria, ALU backend, decoder, disassembler.
- `conformance/`: programmi sintetici senza dipendenze esterne.
- `cmd/retronet-6502/`: runner CLI.
- `docs/`: documentazione utente e architetturale.
