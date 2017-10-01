package bk

import (
	"math/rand"
	"time"
)

const (
	characters = "0123456789"
	charLen    = len(characters)
)

type Secret struct {
	Secret []byte
	Size   uint8
}

func (s Secret) String() string {
	return string(s.Secret)
}

func (s Secret) Guess(guess []byte) Match {
	var extIdx, intIdx uint8
	match := NewMatch(s.Size)
	m := make([]uint8, s.Size)

	for extIdx = 0; extIdx < s.Size; extIdx++ {
		if s.Secret[extIdx] != guess[extIdx] {
			continue
		}

		m[extIdx] = b
		match.b++
	}

	for extIdx = 0; extIdx < s.Size; extIdx++ {
		if m[extIdx] == b {
			continue
		}

		for intIdx = 0; intIdx < s.Size; intIdx++ {
			if m[intIdx] == b {
				continue
			}

			if s.Secret[intIdx] != guess[extIdx] {
				continue
			}

			m[extIdx] = c
			match.c++
			break
		}
	}

	return match
}

type SecretGenerator interface {
	CreateSecret() *Secret
}

type randomSecret struct {
	l uint8
}

func NewRandomSecretGenerator(l uint8) SecretGenerator {
	return randomSecret{l}
}

func (sg randomSecret) CreateSecret() *Secret {
	rand.Seed(time.Now().Unix())

	secret := make([]byte, sg.l)

	for i := range secret {
		secret[i] = characters[rand.Intn(charLen)]
	}

	return &Secret{secret, sg.l}
}

type predefinedSecret struct {
	secret *Secret
}

func NewPredefinedSecretGenerator(secret []byte) SecretGenerator {
	size := uint8(len(secret))
	return predefinedSecret{&Secret{secret, size}}
}

func (sg predefinedSecret) CreateSecret() *Secret {
	return sg.secret
}
