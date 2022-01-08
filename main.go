package main

import (
	"errors"
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

func main() {
}
