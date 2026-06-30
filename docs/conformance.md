# Conformance

La conformance interna non usa ROM storiche. Il pacchetto `conformance` esegue
programmi sintetici auto-verificanti che coprono:

- branch e loop;
- `JSR/RTS` e stack;
- `ADC` in decimal mode;
- bug `JMP ($xxFF)`.

La CLI espone la batteria con:

```bash
go run ./cmd/retronet-6502 -conformance
go run ./cmd/retronet-6502 -conformance -alu native
```

La garanzia piu' importante resta il differenziale tra backend:

- `retronet-hardware/bridge/i6502` confronta la ALU a porte con un riferimento
  Go su tutti gli ingressi di `ADC/SBC`;
- `retronet-6502/cpu` confronta `cpu.Gate` e `cpu.Native` su `ADC/SBC` binari e
  decimali.

Suite esterne per-istruzione possono essere aggiunte in futuro come loader
opzionali, senza vendorizzare dataset non necessari nel repository.
