package chip8

import (
	"github.com/veandco/go-sdl2/sdl"
)

var (
	lo_width  int32 = 64
	lo_height int32 = 32
	hi_width  int32 = 128
	hi_height int32 = 64
)

type Screen struct {
	width, height, scale int32
	window               *sdl.Window
	surface              *sdl.Surface
}

func NewScreen(scale int32) *Screen {
	win := newWindow(lo_width, lo_height, scale)
	surface, err := win.GetSurface()
	if err != nil {
		panic(err)
	}

	screen := Screen{width: lo_width, height: lo_height, scale: scale, window: win, surface: surface}
	return &screen
}

func newWindow(width int32, height int32, scale int32) *sdl.Window {
	win, err := sdl.CreateWindow("CHIP-8 Emulator", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width*scale, height*scale, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	return win
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

func (s *Screen) EnableLowRes() {
	s.window.Destroy()
	win := newWindow(lo_width, lo_height, s.scale)
	surface, err := win.GetSurface()
	if err != nil {
		panic(err)
	}
	s.window = win
	s.surface = surface
	s.width = lo_width
	s.height = lo_height
}

func (s *Screen) EnableHighRes() {
	s.window.Destroy()
	win := newWindow(hi_width, hi_height, s.scale)
	surface, err := win.GetSurface()
	if err != nil {
		panic(err)
	}
	s.window = win
	s.surface = surface
	s.width = hi_width
	s.height = hi_height
}
