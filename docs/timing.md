# Timing

Il core conta cicli **ufficiali aggregati per istruzione**. Non modella mezzi
cicli, segnali elettrici o accessi bus-cycle.

## Regole Implementate

- Ogni opcode ha il ciclo base documentato nella tabella interna `opcodes`.
- Branch: 2 cicli base, +1 se preso, +1 ulteriore se il target e' in un'altra
  pagina.
- Load/ALU/CMP con absolute indexed e indirect indexed: +1 ciclo su page crossing
  dove previsto dall'ISA.
- Store indexed absolute/indirect-indexed usa il ciclo fisso ufficiale, senza
  page-extra dinamico.
- Read-modify-write indexed absolute usa il ciclo fisso ufficiale.

I contatori pubblici sono:

- `InstructionCount`;
- `CycleCount`;
- `LastCycles`.

`Reset()` imposta `LastCycles=7`, ma non incrementa `CycleCount`; il conteggio
inizia dalle istruzioni eseguite dal chiamante.
