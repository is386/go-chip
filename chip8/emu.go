package chip8

import (
	"fmt"
	"io/ioutil"
)

type Emulator struct {
	memory     [4096]uint16
	pc         uint16
	index      uint16
	registers  [16]uint8
	stack      Stack
	soundTimer uint8
	delayTimer uint8
	screen     Screen
}

func NewEmulator(screen Screen) *Emulator {
	emu := Emulator{pc: 0x200}
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

func (emu *Emulator) decode(instr uint16) (uint16, uint16, uint16, uint16, uint16, uint16) {
	op := instr >> 12
	X := (instr & 0x0F00) >> 8
	Y := (instr & 0x00F0) >> 4
	N := instr & 0xF
	NN := instr & 0xFF
	NNN := instr & 0xFFF
	return op, X, Y, N, NN, NNN
}

func (emu *Emulator) Execute() {
	instr := emu.fetch()
	op, _, _, _, _, _ := emu.decode(instr)
	switch op {
	case 0x0:
		fmt.Println("clear screen")
	case 0x1:
		fmt.Println("jump")
	case 0x6:
		fmt.Println("set register")
	case 0x7:
		fmt.Println("add to VX")
	case 0xA:
		fmt.Println("set index")
	case 0xD:
		fmt.Println("display")
	}
}
