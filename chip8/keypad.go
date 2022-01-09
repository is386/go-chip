package chip8

import "github.com/veandco/go-sdl2/sdl"

var (
	keymap = map[sdl.Keycode]uint8{
		sdl.K_x: 0x0,
		sdl.K_1: 0x1,
		sdl.K_2: 0x2,
		sdl.K_3: 0x3,
		sdl.K_q: 0x4,
		sdl.K_w: 0x5,
		sdl.K_e: 0x6,
		sdl.K_a: 0x7,
		sdl.K_s: 0x8,
		sdl.K_d: 0x9,
		sdl.K_z: 0xA,
		sdl.K_c: 0xB,
		sdl.K_4: 0xC,
		sdl.K_r: 0xD,
		sdl.K_f: 0xE,
		sdl.K_v: 0xF,
	}
)

type Keypad struct {
	keys [16]bool
}

func NewKeypad() *Keypad {
	k := Keypad{}
	return &k
}

func (k *Keypad) KeyPressed(key uint8) bool {
	return k.keys[key]
}

func (k *Keypad) KeyEvent(event *sdl.KeyboardEvent) {
	switch event.Type {
	case sdl.KEYUP:
		k.keyUp(event.Keysym.Sym)
	case sdl.KEYDOWN:
		k.keyDown(event.Keysym.Sym)
	}
}

func (k *Keypad) keyDown(key sdl.Keycode) {
	k.keys[keymap[key]] = true
}

func (k *Keypad) keyUp(key sdl.Keycode) {
	k.keys[keymap[key]] = false
}
