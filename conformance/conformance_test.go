package conformance

import (
	"testing"

	"github.com/retronet-labs/retronet-6502/cpu"
)

func TestConformanceGateAndNative(t *testing.T) {
	for _, be := range []struct {
		name string
		alu  cpu.ALUBackend
	}{{"gate", cpu.Gate}, {"native", cpu.Native}} {
		t.Run(be.name, func(t *testing.T) {
			res := Run(be.alu)
			if res.Failed() != 0 {
				t.Fatalf("conformance fallita: %+v", res.Cases)
			}
		})
	}
}
