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

The CHIP-8 had a hexadecimal keypad. The keys are mapped using the table below:

|     CHIP-8      |    Emulator     |
| :-------------: | :-------------: |
| `1` `2` `3` `C` | `1` `2` `3` `4` |
| `4` `5` `6` `D` | `Q` `W` `E` `R` |
| `7` `8` `9` `E` | `A` `S` `D` `F` |
| `A` `0` `B` `F` | `Z` `X` `C` `V` |
