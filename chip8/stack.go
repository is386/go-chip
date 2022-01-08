package chip8

import "errors"

type Stack []uint16

func (s *Stack) Push(val uint16) {
	*s = append(*s, val)
}

func (s *Stack) Pop() (uint16, error) {
	if len(*s) == 0 {
		return 0x0, errors.New("cannot pop from empty stack")
	}
	last := len(*s) - 1
	val := (*s)[last]
	*s = (*s)[:last]
	return val, nil
}
