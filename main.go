package main

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
	FILENAME = "ibm.ch8"
	SCALE    = 10.0
	FONT     = [80]byte{
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

type Stack []byte

func (s *Stack) push(val byte) {
	*s = append(*s, val)
}

func (s *Stack) pop() (byte, error) {
	if len(*s) == 0 {
		return 0x0, errors.New("cannot pop from empty stack")
	}
	last := len(*s) - 1
	val := (*s)[last]
	*s = (*s)[:last]
	return val, nil
}

type Screen struct {
	width, height, scale float64
	window               *pixelgl.Window
	drawer               *imdraw.IMDraw
}

func newScreen(scale float64) *Screen {
	cfg := pixelgl.WindowConfig{
		Title:  "CHIP-8 Emulator",
		Bounds: pixel.R(0, 0, 64*scale, 32*scale),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	screen := Screen{width: 64, height: 32, scale: SCALE, window: win, drawer: imd}
	screen.Clear()
	screen.Update()
	return &screen
}

func (s *Screen) Clear() {
	s.window.Clear(colornames.Black)
}

func (s *Screen) Update() {
	s.window.Update()
}

func (s *Screen) Closed() bool {
	return s.window.Closed()
}

func (s *Screen) DrawPixel(x0 float64, y0 float64, color color.RGBA) {
	x := x0 * s.scale
	y := y0 * s.scale
	s.drawer.Color = color
	s.drawer.Push(pixel.V(x, y))
	s.drawer.Push(pixel.V(x+s.scale, y+s.scale))
	s.drawer.Rectangle(0)
	s.drawer.Draw(s.window)
}

func (s *Screen) GetColor(x float64, y float64) int {
	color := s.window.Color(pixel.V(x*s.scale, y*s.scale))
	fmt.Println(color)
	if color.R == 0 && color.G == 0 && color.B == 0 && color.A == 1 {
		return 0
	}
	return 1
}

func run() {
	screen := newScreen(SCALE)
	for !screen.Closed() {
	}
}

func main() {
	fmt.Println(FILENAME)
	pixelgl.Run(run)
}
