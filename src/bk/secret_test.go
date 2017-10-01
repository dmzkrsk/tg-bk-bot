package bk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	in    string
	guess string
	m     string
}{
	{"1111", "1111", "4B0C"},
	{"1111", "1112", "3B0C"},
	{"1111", "1121", "3B0C"},
	{"1111", "1211", "3B0C"},
	{"1111", "2111", "3B0C"},

	{"1234", "1234", "4B0C"},
	{"1234", "2134", "2B2C"},
	{"1234", "3124", "1B3C"},
	{"1234", "1324", "2B2C"},
	{"1234", "2314", "1B3C"},
	{"1234", "3214", "2B2C"},
	{"1234", "3241", "1B3C"},
	{"1234", "2341", "0B4C"},
	{"1234", "4321", "0B4C"},
	{"1234", "3421", "0B4C"},
	{"1234", "2431", "1B3C"},
	{"1234", "4231", "2B2C"},
	{"1234", "4132", "1B3C"},

	// {"1234", "1432", "4B0C"},
	// {"1234", "3412", "4B0C"},
	// {"1234", "4312", "4B0C"},
	// {"1234", "1342", "4B0C"},
	// {"1234", "3142", "4B0C"},
	// {"1234", "2143", "4B0C"},
	// {"1234", "1243", "4B0C"},
	// {"1234", "4213", "4B0C"},
	// {"1234", "2413", "4B0C"},
	// {"1234", "1423", "4B0C"},
	// {"1234", "4123", "4B0C"},

	{"6436", "4561", "0B2C"},
	{"6436", "5462", "1B1C"},
}

func TestSomething(t *testing.T) {

	for _, test := range tests {
		in := []byte(test.in)
		guess := []byte(test.guess)
		result := NewPredefinedSecretGenerator(in).CreateSecret().Guess(guess).String()
		assert.Equal(t, test.m, result, "%s <> %s", test.in, test.guess)
	}
}
