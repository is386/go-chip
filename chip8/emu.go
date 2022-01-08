package chip8

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
