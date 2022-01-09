package main

import (
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/is386/GoCHIP/chip8"
)

var (
	FILENAME   = "brick.ch8"
	SCALE      = 10.0
	EMU_DELAY  = 10 * time.Millisecond
	TIME_DELAY = 16 * time.Millisecond
)

func decrementTimers(emu *chip8.Emulator, ticker *time.Ticker, stopTicking chan bool) {
	for {
		select {
		case <-stopTicking:
			return
		case <-ticker.C:
			emu.DecrementTimers()
		}
	}
}

func run() {
	screen := chip8.NewScreen(SCALE)
	emu := chip8.NewEmulator(screen)

	ticker := time.NewTicker(TIME_DELAY)
	stopTicking := make(chan bool)
	go decrementTimers(emu, ticker, stopTicking)

	emu.LoadRom(FILENAME)

	for !screen.Closed() {
		emu.Execute()
		screen.Update()
	}

	ticker.Stop()
	stopTicking <- true
}

func main() {
	pixelgl.Run(run)
}
