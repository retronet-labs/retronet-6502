# Flag E Decimal Mode

Il registro P viene impacchettato come `N V 1 B D I Z C`.

## Carry

`C` e' un riporto per `ADC`, ma nelle sottrazioni (`SBC`, `CMP`, `CPX`, `CPY`)
vale **nessun prestito**:

- `C=1`: `A >= operando`;
- `C=0`: c'e' stato prestito.

`SBC` usa quindi `A - value - !C`.

## Overflow

`V` segnala overflow signed nelle operazioni aritmetiche:

- `ADC`: operandi con stesso segno e risultato con segno diverso;
- `SBC`: operandi con segno diverso e risultato con segno diverso da `A`.

`BIT` copia invece il bit 6 dell'operando in `V` e il bit 7 in `N`.

## Decimal Mode

Con `D=1`, `ADC` e `SBC` applicano la correzione BCD. Il modello scelto e'
NMOS 6502:

- il risultato scritto in `A` e' corretto in BCD;
- `C` riflette il carry/prestito decimale;
- `Z`, `N` e `V` derivano dal risultato binario pre-correzione.

Questo comportamento e' documentato nei test del bridge `i6502` e nel
differenziale Gate/Native del core.
