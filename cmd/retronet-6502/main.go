// Comando retronet-6502: esegue binari raw sul core MOS/NMOS 6502.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/retronet-labs/retronet-6502/conformance"
	"github.com/retronet-labs/retronet-6502/cpu"
)

func main() {
	bin := flag.String("bin", "", "file binario raw da caricare")
	load := flag.String("load", "0x8000", "indirizzo di caricamento")
	pc := flag.String("pc", "", "PC iniziale; se vuoto usa il reset vector")
	steps := flag.Int("steps", 100000, "numero massimo di istruzioni")
	trace := flag.Bool("trace", false, "stampa stato e disassembly a ogni istruzione")
	disasm := flag.Int("disasm", 0, "disassembla N istruzioni del binario e termina")
	aluName := flag.String("alu", "gate", "backend ALU: gate oppure native")
	conf := flag.Bool("conformance", false, "esegue la batteria sintetica")
	flag.Parse()

	if *conf {
		runConformance(backendFor(*aluName))
		return
	}
	if *bin == "" {
		flag.Usage()
		os.Exit(2)
	}

	loadAddr := mustAddr(*load)
	if *disasm > 0 {
		listBinary(*bin, loadAddr, *disasm)
		return
	}
	runBinary(*bin, backendFor(*aluName), loadAddr, *pc, *steps, *trace)
}

func backendFor(name string) cpu.ALUBackend {
	if strings.EqualFold(name, "native") {
		return cpu.Native
	}
	return cpu.Gate
}

func runBinary(path string, backend cpu.ALUBackend, loadAddr uint16, pcText string, steps int, trace bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "errore:", err)
		os.Exit(1)
	}
	c := cpu.NewCPU6502WithALU(backend)
	ram := c.Mem.(*cpu.RAM)
	ram.LoadAt(loadAddr, data)
	if pcText == "" {
		ram.Write(0xFFFC, byte(loadAddr))
		ram.Write(0xFFFD, byte(loadAddr>>8))
		c.Reset()
	} else {
		c.PC = mustAddr(pcText)
	}

	executed := 0
	for executed < steps {
		if trace {
			printState(c)
		}
		if err := c.Step(); err != nil {
			fmt.Fprintln(os.Stderr, "stop:", err)
			break
		}
		executed++
	}
	fmt.Printf("eseguite %d istruzioni\n", executed)
	printState(c)
}

func listBinary(path string, loadAddr uint16, count int) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "errore:", err)
		os.Exit(1)
	}
	c := cpu.NewCPU6502()
	c.Mem.(*cpu.RAM).LoadAt(loadAddr, data)
	pc := loadAddr
	for i := 0; i < count; i++ {
		text, n := c.Disassemble(pc)
		fmt.Printf("%04X  %-18s\n", pc, text)
		if n == 0 {
			n = 1
		}
		pc += uint16(n)
	}
}

func printState(c *cpu.CPU6502) {
	text, _ := c.Disassemble(c.PC)
	fmt.Printf("%04X  %-18s A=%02X X=%02X Y=%02X SP=%02X P=%02X CYC=%d\n",
		c.PC, text, c.A, c.X, c.Y, c.SP, c.PackFlags(false), c.CycleCount)
}

func runConformance(backend cpu.ALUBackend) {
	res := conformance.Run(backend)
	for _, c := range res.Cases {
		status := "ok"
		if !c.OK {
			status = "FALLITO: " + c.Detail
		}
		fmt.Printf("[%s] %s\n", status, c.Name)
	}
	fmt.Printf("conformance: %d/%d superati\n", res.Passed(), len(res.Cases))
	if res.Failed() != 0 {
		os.Exit(1)
	}
}

func mustAddr(s string) uint16 {
	v, err := strconv.ParseUint(strings.TrimPrefix(strings.ToLower(strings.TrimSpace(s)), "0x"), 16, 16)
	if err != nil {
		n, decErr := strconv.ParseUint(strings.TrimSpace(s), 10, 16)
		if decErr == nil {
			return uint16(n)
		}
		fmt.Fprintf(os.Stderr, "indirizzo non valido %q\n", s)
		os.Exit(2)
	}
	return uint16(v)
}
