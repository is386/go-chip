package chip8

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/image/colornames"
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
	return &emu
}

func (emu *Emulator) LoadFont(font [80]byte) {
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

func (emu *Emulator) decode(instr uint16) (uint16, uint16, uint16, uint8, uint8, uint16) {
	op := instr >> 12
	X := (instr & 0x0F00) >> 8
	Y := (instr & 0x00F0) >> 4
	N := instr & 0xF
	NN := instr & 0xFF
	NNN := instr & 0xFFF
	return op, X, Y, uint8(N), uint8(NN), NNN
}

func (emu *Emulator) Execute() {
	instr := emu.fetch()
	op, X, Y, N, NN, NNN := emu.decode(instr)
	switch op {
	case 0x0:
		switch NN {
		case 0xE0:
			emu.clearScreen()
		}
	case 0x1:
		emu.jump(NNN)
	case 0x6:
		emu.set_register(X, NN)
	case 0x7:
		emu.add_to_register(X, NN)
	case 0xA:
		emu.set_index(NNN)
	case 0xD:
		emu.display(X, Y, N)
	}
}

func (emu *Emulator) clearScreen() {
	emu.screen.Clear()
}

func (emu *Emulator) jump(NNN uint16) {
	emu.pc = NNN
}

func (emu *Emulator) set_register(X uint16, NN uint8) {
	emu.registers[X] = uint8(NN)
}

func (emu *Emulator) add_to_register(X uint16, NN uint8) {
	temp := emu.registers[X] + NN
	if temp <= 255 {
		emu.registers[X] = temp
	} else {
		emu.registers[X] = temp - 255 - 1
	}
}

func (emu *Emulator) set_index(NNN uint16) {
	emu.index = NNN
}

func (emu *Emulator) display(X uint16, Y uint16, N uint8) {
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
			y := float64(y0 + N - row)

			if x >= emu.screen.width {
				break
			}

			screenPixel := emu.screen.GetColor(x, y)
			if (pixel == 1) && (screenPixel == 1) {
				emu.screen.DrawPixel(x, y, colornames.Black)
			} else if (pixel == 1) && (screenPixel != 1) {
				emu.screen.DrawPixel(x, y, colornames.White)
			}
		}
	}
	emu.screen.Update()
}
