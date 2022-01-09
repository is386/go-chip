package chip8

import (
	"fmt"
	"io/ioutil"
	"math/rand"

	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
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

	keypad = map[uint8]pixelgl.Button{
		0x0: pixelgl.KeyX,
		0x1: pixelgl.Key1,
		0x2: pixelgl.Key2,
		0x3: pixelgl.Key3,
		0x4: pixelgl.KeyQ,
		0x5: pixelgl.KeyW,
		0x6: pixelgl.KeyE,
		0x7: pixelgl.KeyA,
		0x8: pixelgl.KeyS,
		0x9: pixelgl.KeyD,
		0xA: pixelgl.KeyZ,
		0xB: pixelgl.KeyC,
		0xC: pixelgl.Key4,
		0xD: pixelgl.KeyR,
		0xE: pixelgl.KeyF,
		0xF: pixelgl.KeyV,
	}
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
}

func NewEmulator(screen *Screen) *Emulator {
	emu := Emulator{pc: 0x200, screen: screen}
	emu.loadFont()
	return &emu
}

func (emu *Emulator) loadFont() {
	for i := 0; i < len(font); i++ {
		emu.memory[i] = uint16(font[i])
	}
}

func (emu *Emulator) LoadRom(filename string) {
	rom, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(rom); i++ {
		emu.memory[emu.pc+uint16(i)] = uint16(rom[i])
	}
}

func (emu *Emulator) fetch() uint16 {
	byte1 := emu.memory[emu.pc]
	byte2 := emu.memory[emu.pc+1]
	emu.pc += 2
	return (byte1 << 8) | byte2
}

func (emu *Emulator) decode(instr uint16) (uint8, uint8, uint8, uint8, uint8, uint16) {
	op := instr >> 12
	X := (instr & 0x0F00) >> 8
	Y := (instr & 0x00F0) >> 4
	N := instr & 0xF
	NN := instr & 0xFF
	NNN := instr & 0xFFF
	return uint8(op), uint8(X), uint8(Y), uint8(N), uint8(NN), NNN
}

func (emu *Emulator) Execute() {
	instr := emu.fetch()
	op, X, Y, N, NN, NNN := emu.decode(instr)
	switch op {
	case 0x0:
		switch NN {
		case 0xE0:
			emu.clearScreen()
		case 0xEE:
			emu.returnFromSubroutine()
		}
	case 0x1:
		emu.jump(NNN)
	case 0x2:
		emu.callSubroutine(NNN)
	case 0x3:
		if emu.registers[X] == NN {
			emu.skip()
		}
	case 0x4:
		if emu.registers[X] != NN {
			emu.skip()
		}
	case 0x5:
		if emu.registers[X] == emu.registers[Y] {
			emu.skip()
		}
	case 0x9:
		if emu.registers[X] != emu.registers[Y] {
			emu.skip()
		}
	case 0x6:
		emu.setRegister(X, NN)
	case 0x7:
		emu.addToRegister(X, NN)
	case 0x8:
		switch N {
		case 0x0:
			emu.setRegister(X, emu.registers[Y])
		case 0x1:
			emu.logicalOr(X, Y)
		case 0x2:
			emu.logicalAnd(X, Y)
		case 0x3:
			emu.logicalXor(X, Y)
		case 0x4:
			emu.addTwoRegisters(X, Y)
		case 0x5:
			emu.subTwoRegisters(X, Y)
		case 0x6:
			emu.shiftRegisterRight(X)
		case 0x7:
			emu.subnTwoRegisters(X, Y)
		case 0xE:
			emu.shiftRegisterLeft(X)
		}
	case 0xA:
		emu.setIndex(NNN)
	case 0xB:
		emu.jumpWithOffset(NNN)
	case 0xC:
		emu.generateRandom(X, NN)
	case 0xD:
		emu.display(X, Y, N)
	case 0xE:
		switch NN {
		case 0x9E:
			if emu.keyPressed(X) {
				emu.skip()
			}
		case 0xA1:
			if !emu.keyPressed(X) {
				emu.skip()
			}
		}
	case 0xF:
		switch NN {
		case 0x07:
			emu.setRegister(X, emu.delayTimer)
		case 0x15:
			emu.setDelayTimer(X)
		case 0x18:
			emu.setSoundTimer(X)
		case 0x1E:
			emu.addToIndex(X)
		case 0x0A:
			emu.waitForKey(X)
		case 0x29:
			emu.setIndexToFont(X)
		case 0x33:
			emu.storeDecimal(X)
		case 0x55:
			emu.storeRegisters(X)
		case 0x65:
			emu.loadRegisters(X)
		}
	}
}

func (emu *Emulator) clearScreen() {
	emu.screen.Clear()
}

func (emu *Emulator) jump(NNN uint16) {
	emu.pc = NNN
}

func (emu *Emulator) callSubroutine(NNN uint16) {
	emu.stack.Push(emu.pc)
	emu.jump(NNN)
}

func (emu *Emulator) returnFromSubroutine() {
	addr, err := emu.stack.Pop()
	if err != nil {
		panic(err)
	}
	emu.pc = addr
}

func (emu *Emulator) skip() {
	emu.pc += 2
}

func (emu *Emulator) setRegister(X uint8, NN uint8) {
	emu.registers[X] = NN
}

func (emu *Emulator) addToRegister(X uint8, NN uint8) {
	emu.registers[X] = emu.registers[X] + NN
}

func (emu *Emulator) logicalOr(X uint8, Y uint8) {
	emu.registers[X] |= emu.registers[Y]
}

func (emu *Emulator) logicalAnd(X uint8, Y uint8) {
	emu.registers[X] &= emu.registers[Y]
}

func (emu *Emulator) logicalXor(X uint8, Y uint8) {
	emu.registers[X] ^= emu.registers[Y]
}

func (emu *Emulator) addTwoRegisters(X uint8, Y uint8) {
	temp := uint16(emu.registers[X] + emu.registers[Y])
	if temp <= 255 {
		emu.registers[X] = uint8(temp)
		emu.registers[0xF] = 0
	} else {
		emu.registers[X] = uint8(temp - 256)
		emu.registers[0xF] = 1
	}
}

func (emu *Emulator) subTwoRegisters(X uint8, Y uint8) {
	if emu.registers[X] > emu.registers[Y] {
		emu.registers[X] -= emu.registers[Y]
		emu.registers[0xF] = 1
	} else {
		emu.registers[0xF] = 0
		emu.registers[X] -= emu.registers[Y] + 255 + 1
	}
}

func (emu *Emulator) subnTwoRegisters(X uint8, Y uint8) {
	if emu.registers[Y] > emu.registers[X] {
		emu.registers[X] = emu.registers[Y] - emu.registers[X]
		emu.registers[0xF] = 1
	} else {
		emu.registers[0xF] = 0
		emu.registers[X] = emu.registers[Y] - emu.registers[X] + 255 + 1
	}
}

func (emu *Emulator) shiftRegisterRight(X uint8) {
	emu.registers[0xF] = emu.registers[X] & 0x1
	emu.registers[X] >>= 1
}

func (emu *Emulator) shiftRegisterLeft(X uint8) {
	emu.registers[0xF] = (emu.registers[X] & 0x80) >> 7
	emu.registers[X] <<= 1
}

func (emu *Emulator) setIndex(NNN uint16) {
	emu.index = NNN
}

func (emu *Emulator) jumpWithOffset(NNN uint16) {
	emu.pc = uint16(emu.registers[0]) + NNN
}

func (emu *Emulator) generateRandom(X uint8, NN uint8) {
	emu.registers[X] = uint8(rand.Intn(255)) & NN
}

func (emu *Emulator) display(X uint8, Y uint8, N uint8) {
	x0 := emu.registers[X] % uint8(emu.screen.width)
	y0 := emu.registers[Y] % uint8(emu.screen.height)
	emu.registers[0xF] = 0

	var row uint8
	for row = 0; row < N; row++ {
		if y0+row >= uint8(emu.screen.height) {
			break
		}

		sprite := emu.memory[emu.index+uint16(row)]
		spriteBin := fmt.Sprintf("%08b", sprite)

		var col uint8
		for col = 0; col < 8; col++ {
			pixel := int(spriteBin[col] - '0')
			x := float64(x0 + col)
			y := float64(uint8(emu.screen.height) - y0 - row - 1)

			if x >= emu.screen.width {
				break
			}

			screenPixel := emu.screen.GetColor(x, y)
			if (pixel == 1) && (screenPixel == 1) {
				emu.screen.DrawPixel(x, y, colornames.Black)
				emu.registers[0xF] = 1
			} else if (pixel == 1) && (screenPixel != 1) {
				emu.screen.DrawPixel(x, y, colornames.White)
			}
		}
	}
}

func (emu *Emulator) keyPressed(X uint8) bool {
	return emu.screen.window.Pressed(keypad[emu.registers[X]])
}

func (emu *Emulator) setDelayTimer(X uint8) {
	emu.delayTimer = emu.registers[X]
}

func (emu *Emulator) setSoundTimer(X uint8) {
	emu.soundTimer = emu.registers[X]
}

func (emu *Emulator) addToIndex(X uint8) {
	emu.index += uint16(emu.registers[X])
}

func (emu *Emulator) waitForKey(X uint8) {
	for {
		emu.screen.window.UpdateInputWait(0)
		for keyByte, key := range keypad {
			if emu.screen.window.JustPressed(key) {
				emu.registers[X] = keyByte
			}
		}
	}
}

func (emu *Emulator) setIndexToFont(X uint8) {
	emu.index = 5 * uint16(emu.registers[X])
}

func (emu *Emulator) storeDecimal(X uint8) {
	emu.memory[emu.index] = uint16(emu.registers[X] / 100)
	emu.memory[emu.index+1] = uint16(emu.registers[X] % 100 / 10)
	emu.memory[emu.index+2] = uint16(emu.registers[X] % 10)
}

func (emu *Emulator) storeRegisters(X uint8) {
	var i uint8
	for i = 0; i < X+1; i++ {
		idx := uint16(i)
		emu.memory[emu.index+idx] = uint16(emu.registers[idx])
	}
}

func (emu *Emulator) loadRegisters(X uint8) {
	var i uint8
	for i = 0; i < X+1; i++ {
		idx := uint16(i)
		emu.registers[idx] = uint8(emu.memory[emu.index+idx])
	}
}

func (emu *Emulator) DecrementTimers() {
	if emu.delayTimer != 0 {
		emu.delayTimer -= 1
	}
	if emu.soundTimer != 0 {
		emu.soundTimer -= 1
	}
}
