package random

import (
	"crypto/rand"
	"errors"
)

// Character-class flags
const (
	Lower   uint8 = 1 << iota // a-z
	Upper                     // A-Z
	Digits                    // 0-9
	Symbols                   // !@#$%^&*()-_=+[]{}<>,.?/:;|~
)

// pools must stay in the same order as the iota above
var pools = []string{
	"abcdefghijklmnopqrstuvwxyz",
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"0123456789",
	"!@#$%^&*()-_=+[]{}<>,.?/:;|~",
}

// RandomString returns a cryptographically-secure random string of length n
// built from the character classes selected by flags (OR-ed together).
// Passing 0 for flags is an error.
func RandomString(n int, flags uint8) (string, error) {
	if n <= 0 {
		return "", errors.New("length must be positive")
	}
	if flags == 0 {
		return "", errors.New("at least one character class must be specified")
	}

	// Build the requested alphabet
	var alphabet string
	for i, p := range pools {
		if flags&(1<<i) != 0 {
			alphabet += p
		}
	}
	if alphabet == "" {
		return "", errors.New("empty alphabet")
	}

	// Generate
	letterCount := byte(len(alphabet))
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i, v := range b {
		b[i] = alphabet[v%letterCount]
	}
	return string(b), nil
}
