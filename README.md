# GoCHIP

This is a simple CHIP-8 emulator written in Go. Check out this blogpost to learn more about CHIP-8 emulation: https://tobiasvl.github.io/blog/write-a-chip-8-emulator/. My implementation does not emulate sound.

## Usage

`go run main.go <rom_file>`

The `<rom_file>` parameter is just a CHIP-8 ROM. I have provided a few in this repository in `/roms/`.

## Dependencies

- `go 1.5`

### Go Dependencies

- `github.com/veandco/go-sdl2/sdl`
