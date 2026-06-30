# Architettura

`retronet-6502` separa tre responsabilita':

- `cpu.CPU6502`: stato registri/flag, fetch-decode-execute, stack, interrupt e
  conteggio cicli;
- `cpu.Bus`: interfaccia memoria a 16 bit (`Read`/`Write`);
- `cpu.ALUBackend`: motore aritmetico-logico, scegliibile tra `Gate` e `Native`.

## Stato CPU

La CPU espone registri `A`, `X`, `Y`, `SP`, `PC` e flag booleani:

| Flag | Nome | Note |
|------|------|------|
| `N` | Negative | copia bit 7 dei risultati che aggiornano `N/Z` |
| `V` | Overflow | overflow signed per `ADC/SBC`, bit 6 dell'operando per `BIT` |
| `D` | Decimal | abilita correzione BCD per `ADC/SBC` |
| `I` | Interrupt disable | maschera `IRQ`, non `NMI` |
| `Z` | Zero | risultato nullo |
| `C` | Carry | riporto; nelle sottrazioni significa nessun prestito |

Il bit `B` non e' stato interno reale: viene sintetizzato da `PackFlags` quando
`PHP` o `BRK` salvano il registro P sullo stack. Il bit 5 e' sempre salvato a 1.

## Reset E Interrupt

`Reset()` legge il vettore `$FFFC/$FFFD`, imposta `SP=$FD`, abilita `I` e azzera
`D`. La memoria non viene cancellata.

`IRQ()` consegna un interrupt solo se `I=0`; `NMI()` e' sempre accettato.
`BRK` salta al vettore IRQ/BRK ma salva il flag `B=1`; IRQ/NMI salvano `B=0`.
`RTI` ripristina P e PC dallo stack.

## ALU Gate/Native

`Gate` delega a `github.com/retronet-labs/retronet-hardware/bridge/i6502`, che
adatta la ALU a porte di RetroNet Logic alla semantica del 6502. `Native`
riproduce la stessa semantica con operatori Go. Il test
`TestGateVsNativeALUDifferential` confronta i due backend su tutti gli ingressi
di `ADC/SBC`, inclusa la modalita' decimale.

Il backend e' configurazione, non stato:

```go
c := cpu.NewCPU6502WithALU(cpu.Native)
c.SetALU(cpu.Gate)
```
