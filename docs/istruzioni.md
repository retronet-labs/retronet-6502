# Istruzioni E Addressing Mode

Il core implementa gli opcode **documentati** del 6502 NMOS. Gli opcode illegali
non vengono reinterpretati: `Step` restituisce `*cpu.IllegalOpcodeError`.

## Famiglie Implementate

- Load/store: `LDA`, `LDX`, `LDY`, `STA`, `STX`, `STY`.
- ALU/logica: `ADC`, `SBC`, `AND`, `ORA`, `EOR`, `CMP`, `CPX`, `CPY`, `BIT`.
- Read-modify-write: `ASL`, `LSR`, `ROL`, `ROR`, `INC`, `DEC`.
- Registri: `TAX`, `TAY`, `TXA`, `TYA`, `TSX`, `TXS`, `INX`, `INY`, `DEX`, `DEY`.
- Controllo: `JMP`, `JSR`, `RTS`, `RTI`, `BRK`, `NOP`.
- Branch: `BPL`, `BMI`, `BVC`, `BVS`, `BCC`, `BCS`, `BNE`, `BEQ`.
- Stack: `PHA`, `PLA`, `PHP`, `PLP`.
- Flag: `CLC`, `SEC`, `CLI`, `SEI`, `CLV`, `CLD`, `SED`.

## Addressing Mode

Sono supportate tutte le forme documentate:

| Modo | Esempio | Note |
|------|---------|------|
| implied | `CLC` | nessun operando |
| accumulator | `ASL A` | opera su `A` |
| immediate | `LDA #$10` | operando nel byte successivo |
| zero page | `LDA $20` | indirizzo `$00xx` |
| zero page indexed | `LDA $20,X` | indice con wrap a 8 bit |
| absolute | `LDA $2000` | indirizzo little-endian |
| absolute indexed | `LDA $2000,Y` | ciclo extra se attraversa pagina, dove previsto |
| indexed indirect | `LDA ($20,X)` | puntatore in zero page dopo somma con `X` |
| indirect indexed | `LDA ($20),Y` | puntatore zero page, poi somma con `Y` |
| indirect | `JMP ($1234)` | solo `JMP`; include bug `$xxFF` |
| relative | `BNE loop` | offset signed a 8 bit |

## Bug Storico `JMP ($xxFF)`

Il 6502 NMOS legge il byte alto del target da `$xx00`, non da `$(xx+1)00`, quando
il puntatore indiretto termina a fine pagina. Il core replica questo comportamento
in `read16Bug`.
