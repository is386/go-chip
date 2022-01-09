package chip8

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Screen struct {
	width, height, scale int32
	window               *sdl.Window
	surface              *sdl.Surface
}

func NewScreen(scale int32) *Screen {
	win, err := sdl.CreateWindow("CHIP-8 Emulator", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		64*scale, 32*scale, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	surface, err := win.GetSurface()
	if err != nil {
		panic(err)
	}

	screen := Screen{width: 64, height: 32, scale: scale, window: win, surface: surface}
	return &screen
}

func (s *Screen) Clear() {
	s.surface.FillRect(nil, 0)
}

func (s *Screen) DrawPixel(x0 int32, y0 int32, color uint32) {
	x := x0 * s.scale
	y := y0 * s.scale
	pixel := sdl.Rect{X: x, Y: y, W: s.scale, H: s.scale}
	s.surface.FillRect(&pixel, color)
	s.window.UpdateSurface()
}

func (s *Screen) GetColor(x int32, y int32) int {
	color := s.surface.At(int(x*s.scale), int(y*s.scale))
	r, g, b, _ := color.RGBA()
	if r+g+b == 0 {
		return 0
	}
	return 1
}
