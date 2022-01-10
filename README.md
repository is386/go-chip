# GoCHIP

This is a simple CHIP-8 emulator written in Go. Check out this blogpost to learn more about CHIP-8 emulation: https://tobiasvl.github.io/blog/write-a-chip-8-emulator/. My implementation does not emulate sound.

## Usage

`go run main.go <rom_file>`

The `<rom_file>` parameter is just a CHIP-8 ROM. I have provided a few in this repository in `/roms/`.

## Dependencies

- `go 1.15`

### Go Dependencies

- `github.com/veandco/go-sdl2/sdl`

## Keyboard

The controls will vary from game to game. These are the keys used by the emulators, it ressembles a hexadecimal keypad:

```
Hexadecimal Keypad:
1 2 3 C
4 5 6 D
7 8 9 E
A 0 B F

Emulated Keypad:
1 2 3 4
Q W E R
A S D F
Z X C V
```
