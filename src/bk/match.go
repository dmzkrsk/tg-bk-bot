package bk

import (
	"fmt"
)

const (
	c uint8 = 1
	b       = 2
)

type Match struct {
	b uint8
	c uint8

	l uint8
}

func NewMatch(l uint8) Match {
	return Match{l: l}
}

func (m Match) IsWin() bool {
	return m.b == m.l && m.c == 0
}

func (m Match) String() string {
	return fmt.Sprintf("%dB%dC", m.b, m.c)
}
