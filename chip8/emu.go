package chip8

import (
	"fmt"
	"io/ioutil"
	"math/rand"
)

var (
	font = [80]byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0,
		0x20, 0x60, 0x20, 0x20, 0x70,
		0xF0, 0x10, 0xF0, 0x80, 0xF0,
		0xF0, 0x10, 0xF0, 0x10, 0xF0,
		0x90, 0x90, 0xF0, 0x10, 0x10,
		0xF0, 0x80, 0xF0, 0x10, 0xF0,
		0xF0, 0x80, 0xF0, 0x90, 0xF0,
		0xF0, 0x10, 0x20, 0x40, 0x40,
		0xF0, 0x90, 0xF0, 0x90, 0xF0,
		0xF0, 0x90, 0xF0, 0x10, 0xF0,
		0xF0, 0x90, 0xF0, 0x90, 0x90,
		0xE0, 0x90, 0xE0, 0x90, 0xE0,
		0xF0, 0x80, 0x80, 0x80, 0xF0,
		0xE0, 0x90, 0x90, 0x90, 0xE0,
		0xF0, 0x80, 0xF0, 0x80, 0xF0,
		0xF0, 0x80, 0xF0, 0x80, 0x80}
)

type Emulator struct {
	memory     [4096]uint16
	pc         uint16
	index      uint16
	registers  [16]uint8
	stack      Stack
	soundTimer uint8
	delayTimer uint8
	screen     *Screen
	keypad     *Keypad
}

func NewEmulator(screen *Screen, keypad *Keypad) *Emulator {
	e := Emulator{pc: 0x200, screen: screen, keypad: keypad}
	e.loadFont()
	return &e
}

func (e *Emulator) loadFont() {
	for i := 0; i < len(font); i++ {
		e.memory[i] = uint16(font[i])
	}
}

func (e *Emulator) LoadRom(filename string) {
	rom, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(rom); i++ {
		e.memory[e.pc+uint16(i)] = uint16(rom[i])
	}
}

func (e *Emulator) fetch() uint16 {
	byte1 := e.memory[e.pc]
	byte2 := e.memory[e.pc+1]
	e.pc += 2
	return (byte1 << 8) | byte2
}

func (e *Emulator) decode(instr uint16) (uint8, uint8, uint8, uint8, uint8, uint16) {
	op := instr >> 12
	X := (instr & 0x0F00) >> 8
	Y := (instr & 0x00F0) >> 4
	N := instr & 0xF
	NN := instr & 0xFF
	NNN := instr & 0xFFF
	return uint8(op), uint8(X), uint8(Y), uint8(N), uint8(NN), NNN
}

func (e *Emulator) Execute() {
	instr := e.fetch()
	op, X, Y, N, NN, NNN := e.decode(instr)
	switch op {
	case 0x0:
		switch NN {
		case 0xE0:
			e.clearScreen()
		case 0xEE:
			e.returnFromSubroutine()
		}
	case 0x1:
		e.jump(NNN)
	case 0x2:
		e.callSubroutine(NNN)
	case 0x3:
		if e.registers[X] == NN {
			e.skip()
		}
	case 0x4:
		if e.registers[X] != NN {
			e.skip()
		}
	case 0x5:
		if e.registers[X] == e.registers[Y] {
			e.skip()
		}
	case 0x9:
		if e.registers[X] != e.registers[Y] {
			e.skip()
		}
	case 0x6:
		e.setRegister(X, NN)
	case 0x7:
		e.addToRegister(X, NN)
	case 0x8:
		switch N {
		case 0x0:
			e.setRegister(X, e.registers[Y])
		case 0x1:
			e.logicalOr(X, Y)
		case 0x2:
			e.logicalAnd(X, Y)
		case 0x3:
			e.logicalXor(X, Y)
		case 0x4:
			e.addTwoRegisters(X, Y)
		case 0x5:
			e.subTwoRegisters(X, Y)
		case 0x6:
			e.shiftRegisterRight(X)
		case 0x7:
			e.subnTwoRegisters(X, Y)
		case 0xE:
			e.shiftRegisterLeft(X)
		}
	case 0xA:
		e.setIndex(NNN)
	case 0xB:
		e.jumpWithOffset(NNN)
	case 0xC:
		e.generateRandom(X, NN)
	case 0xD:
		e.display(X, Y, N)
	case 0xE:
		switch NN {
		case 0x9E:
			if e.keyPressed(X) {
				e.skip()
			}
		case 0xA1:
			if !e.keyPressed(X) {
				e.skip()
			}
		}
	case 0xF:
		switch NN {
		case 0x07:
			e.setRegister(X, e.delayTimer)
		case 0x15:
			e.setDelayTimer(X)
		case 0x18:
			e.setSoundTimer(X)
		case 0x1E:
			e.addToIndex(X)
		case 0x0A:
			e.waitForKey(X)
		case 0x29:
			e.setIndexToFont(X)
		case 0x33:
			e.storeDecimal(X)
		case 0x55:
			e.storeRegisters(X)
		case 0x65:
			e.loadRegisters(X)
		}
	}
}

func (e *Emulator) clearScreen() {
	e.screen.Clear()
}

func (e *Emulator) jump(NNN uint16) {
	e.pc = NNN
}

func (e *Emulator) callSubroutine(NNN uint16) {
	e.stack.Push(e.pc)
	e.jump(NNN)
}

func (e *Emulator) returnFromSubroutine() {
	addr, err := e.stack.Pop()
	if err != nil {
		panic(err)
	}
	e.pc = addr
}

func (e *Emulator) skip() {
	e.pc += 2
}

func (e *Emulator) setRegister(X uint8, NN uint8) {
	e.registers[X] = NN
}

func (e *Emulator) addToRegister(X uint8, NN uint8) {
	e.registers[X] = e.registers[X] + NN
}

func (e *Emulator) logicalOr(X uint8, Y uint8) {
	e.registers[X] |= e.registers[Y]
}

func (e *Emulator) logicalAnd(X uint8, Y uint8) {
	e.registers[X] &= e.registers[Y]
}

func (e *Emulator) logicalXor(X uint8, Y uint8) {
	e.registers[X] ^= e.registers[Y]
}

func (e *Emulator) addTwoRegisters(X uint8, Y uint8) {
	temp := uint16(e.registers[X] + e.registers[Y])
	if temp <= 255 {
		e.registers[X] = uint8(temp)
		e.registers[0xF] = 0
	} else {
		e.registers[X] = uint8(temp - 256)
		e.registers[0xF] = 1
	}
}

func (e *Emulator) subTwoRegisters(X uint8, Y uint8) {
	if e.registers[X] > e.registers[Y] {
		e.registers[X] -= e.registers[Y]
		e.registers[0xF] = 1
	} else {
		e.registers[0xF] = 0
		e.registers[X] -= e.registers[Y] + 255 + 1
	}
}

func (e *Emulator) subnTwoRegisters(X uint8, Y uint8) {
	if e.registers[Y] > e.registers[X] {
		e.registers[X] = e.registers[Y] - e.registers[X]
		e.registers[0xF] = 1
	} else {
		e.registers[0xF] = 0
		e.registers[X] = e.registers[Y] - e.registers[X] + 255 + 1
	}
}

func (e *Emulator) shiftRegisterRight(X uint8) {
	e.registers[0xF] = e.registers[X] & 0x1
	e.registers[X] >>= 1
}

func (e *Emulator) shiftRegisterLeft(X uint8) {
	e.registers[0xF] = (e.registers[X] & 0x80) >> 7
	e.registers[X] <<= 1
}

func (e *Emulator) setIndex(NNN uint16) {
	e.index = NNN
}

func (e *Emulator) jumpWithOffset(NNN uint16) {
	e.pc = uint16(e.registers[0]) + NNN
}

func (e *Emulator) generateRandom(X uint8, NN uint8) {
	e.registers[X] = uint8(rand.Intn(255)) & NN
}

func (e *Emulator) display(X uint8, Y uint8, N uint8) {
	x0 := e.registers[X] % uint8(e.screen.width)
	y0 := e.registers[Y] % uint8(e.screen.height)
	e.registers[0xF] = 0

	var row uint8
	for row = 0; row < N; row++ {
		if y0+row >= uint8(e.screen.height) {
			break
		}

		sprite := e.memory[e.index+uint16(row)]
		spriteBin := fmt.Sprintf("%08b", sprite)

		var col uint8
		for col = 0; col < 8; col++ {
			pixel := int(spriteBin[col] - '0')
			x := int32(x0 + col)
			y := int32(y0 + row)

			if x >= e.screen.width {
				break
			}

			screenPixel := e.screen.GetColor(x, y)
			if (pixel == 1) && (screenPixel == 1) {
				e.screen.DrawPixel(x, y, 0)
				e.registers[0xF] = 1
			} else if (pixel == 1) && (screenPixel != 1) {
				e.screen.DrawPixel(x, y, 0xffffff)
			}
		}
	}
}

func (e *Emulator) keyPressed(X uint8) bool {
	return e.keypad.KeyPressed(e.registers[X])
}

func (e *Emulator) setDelayTimer(X uint8) {
	e.delayTimer = e.registers[X]
}

func (e *Emulator) setSoundTimer(X uint8) {
	e.soundTimer = e.registers[X]
}

func (e *Emulator) addToIndex(X uint8) {
	e.index += uint16(e.registers[X])
}

func (e *Emulator) waitForKey(X uint8) {
	e.registers[X] = e.keypad.WaitForKey()
}

func (e *Emulator) setIndexToFont(X uint8) {
	e.index = 5 * uint16(e.registers[X])
}

func (e *Emulator) storeDecimal(X uint8) {
	e.memory[e.index] = uint16(e.registers[X] / 100)
	e.memory[e.index+1] = uint16(e.registers[X] % 100 / 10)
	e.memory[e.index+2] = uint16(e.registers[X] % 10)
}

func (e *Emulator) storeRegisters(X uint8) {
	var i uint8
	for i = 0; i < X+1; i++ {
		idx := uint16(i)
		e.memory[e.index+idx] = uint16(e.registers[idx])
	}
}

func (e *Emulator) loadRegisters(X uint8) {
	var i uint8
	for i = 0; i < X+1; i++ {
		idx := uint16(i)
		e.registers[idx] = uint8(e.memory[e.index+idx])
	}
}

func (e *Emulator) DecrementTimers() {
	if e.delayTimer != 0 {
		e.delayTimer -= 1
	}
	if e.soundTimer != 0 {
		e.soundTimer -= 1
	}
}
