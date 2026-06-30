package cpu

import "github.com/retronet-labs/retronet-hardware/bridge/i6502"

// Step esegue una singola istruzione documentata del 6502.
func (c *CPU6502) Step() error {
	if c.Mem == nil {
		return ErrNilBus
	}
	pcBefore := c.PC
	code := c.fetch8()
	def := opcodes[code]
	if def.op == opIllegal {
		return &IllegalOpcodeError{PC: pcBefore, Opcode: code}
	}

	cycles := int(def.cycles)
	switch def.op {
	case opADC, opAND, opEOR, opORA, opSBC, opCMP, opCPX, opCPY, opLDA, opLDX, opLDY, opBIT:
		value, crossed := c.readOperand(def.mode)
		if def.pageExtra && crossed {
			cycles++
		}
		c.executeRead(def.op, value)

	case opSTA, opSTX, opSTY:
		addr, _ := c.operandAddr(def.mode)
		switch def.op {
		case opSTA:
			c.write(addr, c.A)
		case opSTX:
			c.write(addr, c.X)
		default:
			c.write(addr, c.Y)
		}

	case opASL, opLSR, opROL, opROR, opINC, opDEC:
		c.executeModify(def)

	case opBPL, opBMI, opBVC, opBVS, opBCC, opBCS, opBNE, opBEQ:
		offset := int8(c.fetch8())
		if c.branchTaken(def.op) {
			old := c.PC
			c.PC = uint16(int(c.PC) + int(offset))
			cycles++
			if !samePage(old, c.PC) {
				cycles++
			}
		}

	case opJMP:
		if def.mode == modeIndirect {
			ptr := c.fetch16()
			c.PC = c.read16Bug(ptr)
		} else {
			c.PC = c.fetch16()
		}
	case opJSR:
		target := c.fetch16()
		c.push16(c.PC - 1)
		c.PC = target
	case opRTS:
		c.PC = c.pop16() + 1
	case opRTI:
		c.SetFlagsByte(c.pop())
		c.PC = c.pop16()
	case opBRK:
		c.PC++ // byte firma ignorato
		c.interrupt(vectorIRQ, true, cycles)
		c.InstructionCount++
		return nil

	case opPHA:
		c.push(c.A)
	case opPHP:
		c.push(c.PackFlags(true))
	case opPLA:
		c.A = c.pop()
		c.setNZ(c.A)
	case opPLP:
		c.SetFlagsByte(c.pop())

	case opCLC:
		c.Carry = false
	case opSEC:
		c.Carry = true
	case opCLI:
		c.InterruptDisable = false
	case opSEI:
		c.InterruptDisable = true
	case opCLV:
		c.Overflow = false
	case opCLD:
		c.Decimal = false
	case opSED:
		c.Decimal = true

	case opTAX:
		c.X = c.A
		c.setNZ(c.X)
	case opTAY:
		c.Y = c.A
		c.setNZ(c.Y)
	case opTSX:
		c.X = c.SP
		c.setNZ(c.X)
	case opTXA:
		c.A = c.X
		c.setNZ(c.A)
	case opTXS:
		c.SP = c.X
	case opTYA:
		c.A = c.Y
		c.setNZ(c.A)
	case opINX:
		c.X = c.incReg(c.X)
	case opINY:
		c.Y = c.incReg(c.Y)
	case opDEX:
		c.X = c.decReg(c.X)
	case opDEY:
		c.Y = c.decReg(c.Y)
	case opNOP:
		// niente
	}

	c.LastCycles = cycles
	c.CycleCount += uint64(cycles)
	c.InstructionCount++
	return nil
}

// Run esegue fino a maxSteps istruzioni o fino al primo errore.
func (c *CPU6502) Run(maxSteps int) (int, error) {
	for i := 0; i < maxSteps; i++ {
		if err := c.Step(); err != nil {
			return i, err
		}
	}
	return maxSteps, nil
}

func (c *CPU6502) executeRead(op operation, value byte) {
	switch op {
	case opADC:
		result, flags := c.backend().ADC(c.A, value, c.Carry, c.Decimal)
		c.A = result
		c.applyALUFlags(flags)
	case opSBC:
		result, flags := c.backend().SBC(c.A, value, c.Carry, c.Decimal)
		c.A = result
		c.applyALUFlags(flags)
	case opAND:
		result, flags := c.backend().Logic(i6502.OpAND, c.A, value)
		c.A = result
		c.applyNZ(flags)
	case opORA:
		result, flags := c.backend().Logic(i6502.OpORA, c.A, value)
		c.A = result
		c.applyNZ(flags)
	case opEOR:
		result, flags := c.backend().Logic(i6502.OpEOR, c.A, value)
		c.A = result
		c.applyNZ(flags)
	case opCMP:
		_, flags := c.backend().Compare(c.A, value)
		c.Carry = flags.Carry
		c.applyNZ(flags)
	case opCPX:
		_, flags := c.backend().Compare(c.X, value)
		c.Carry = flags.Carry
		c.applyNZ(flags)
	case opCPY:
		_, flags := c.backend().Compare(c.Y, value)
		c.Carry = flags.Carry
		c.applyNZ(flags)
	case opLDA:
		c.A = value
		c.setNZ(c.A)
	case opLDX:
		c.X = value
		c.setNZ(c.X)
	case opLDY:
		c.Y = value
		c.setNZ(c.Y)
	case opBIT:
		flags := c.backend().BIT(c.A, value)
		c.Zero = flags.Zero
		c.Negative = flags.Negative
		c.Overflow = flags.Overflow
	}
}

func (c *CPU6502) executeModify(def opcodeDef) {
	if def.mode == modeAccumulator {
		c.A = c.modify(def.op, c.A)
		return
	}
	addr, _ := c.operandAddr(def.mode)
	c.write(addr, c.modify(def.op, c.read(addr)))
}

func (c *CPU6502) modify(op operation, value byte) byte {
	var out byte
	var flags i6502.Flags
	switch op {
	case opASL:
		out, flags = c.backend().ShiftLeft(value)
	case opLSR:
		out, flags = c.backend().ShiftRight(value)
	case opROL:
		out, flags = c.backend().RotateLeft(value, c.Carry)
	case opROR:
		out, flags = c.backend().RotateRight(value, c.Carry)
	case opINC:
		out, flags = c.backend().Increment(value)
	case opDEC:
		out, flags = c.backend().Decrement(value)
	}
	if op == opASL || op == opLSR || op == opROL || op == opROR {
		c.Carry = flags.Carry
	}
	c.applyNZ(flags)
	return out
}

func (c *CPU6502) incReg(v byte) byte {
	out, flags := c.backend().Increment(v)
	c.applyNZ(flags)
	return out
}

func (c *CPU6502) decReg(v byte) byte {
	out, flags := c.backend().Decrement(v)
	c.applyNZ(flags)
	return out
}

func (c *CPU6502) branchTaken(op operation) bool {
	switch op {
	case opBPL:
		return !c.Negative
	case opBMI:
		return c.Negative
	case opBVC:
		return !c.Overflow
	case opBVS:
		return c.Overflow
	case opBCC:
		return !c.Carry
	case opBCS:
		return c.Carry
	case opBNE:
		return !c.Zero
	default:
		return c.Zero
	}
}

func (c *CPU6502) readOperand(mode addrMode) (byte, bool) {
	addr, crossed := c.operandAddr(mode)
	return c.read(addr), crossed
}

func (c *CPU6502) operandAddr(mode addrMode) (uint16, bool) {
	switch mode {
	case modeImmediate:
		addr := c.PC
		c.PC++
		return addr, false
	case modeZeroPage:
		return uint16(c.fetch8()), false
	case modeZeroPageX:
		return uint16(byte(c.fetch8() + c.X)), false
	case modeZeroPageY:
		return uint16(byte(c.fetch8() + c.Y)), false
	case modeAbsolute:
		return c.fetch16(), false
	case modeAbsoluteX:
		base := c.fetch16()
		addr := base + uint16(c.X)
		return addr, !samePage(base, addr)
	case modeAbsoluteY:
		base := c.fetch16()
		addr := base + uint16(c.Y)
		return addr, !samePage(base, addr)
	case modeIndexedIndirect:
		zp := byte(c.fetch8() + c.X)
		lo := uint16(c.read(uint16(zp)))
		hi := uint16(c.read(uint16(byte(zp + 1))))
		return lo | hi<<8, false
	case modeIndirectIndexed:
		zp := c.fetch8()
		base := uint16(c.read(uint16(zp))) | uint16(c.read(uint16(byte(zp+1))))<<8
		addr := base + uint16(c.Y)
		return addr, !samePage(base, addr)
	default:
		return 0, false
	}
}
