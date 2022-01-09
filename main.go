package main

import (
	"time"

	"github.com/is386/GoCHIP/chip8"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	FILENAME           = "roms/asd.ch8"
	SCALE        int32 = 10
	EMU_DELAY          = 2 * time.Millisecond
	TICKER_DELAY       = 16 * time.Millisecond
)

func emulatorTicker(emu *chip8.Emulator, ticker *time.Ticker, stopTicking chan bool) {
	for {
		select {
		case <-stopTicking:
			return
		case <-ticker.C:
			emu.DecrementTimers()
		}
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	screen := chip8.NewScreen(SCALE)
	keypad := chip8.NewKeypad()
	emu := chip8.NewEmulator(screen, keypad)
	emu.LoadRom(FILENAME)

	ticker := time.NewTicker(TICKER_DELAY)
	stopTicking := make(chan bool)
	go emulatorTicker(emu, ticker, stopTicking)

	running := true
	for running {
		time.Sleep(EMU_DELAY)
		emu.Execute()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				keypad.KeyEvent(e)
			}
		}
	}

	ticker.Stop()
	stopTicking <- true
}
