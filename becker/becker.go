// Package becker computes the anti-theft unlock code for Becker car radios
// (the 6-button Mercedes "Special" family and its 4-/8-button siblings) from
// the last four digits of the radio's serial number.
//
// The algorithm was recovered by reverse-engineering the original Becker code
// generator and verified to reproduce the real program's output for all 10000
// possible inputs across every supported variant.
package becker

import (
	"errors"
	"fmt"
)

// Variant selects the Becker keypad family. Its numeric value is the modulus
// (the highest allowed digit), which equals the number of preset buttons used
// to key the code in on the radio.
type Variant int

const (
	// FourButton: 5-digit code, each digit 1-4.
	FourButton Variant = 4
	// SixButton: 5-digit code, each digit 1-6. The Mercedes Becker "Special"
	// family: BE2210, BE1650, BE1350, BE2211.
	SixButton Variant = 6
	// EightButton: 5-digit code, each digit 1-8.
	EightButton Variant = 8
)

// String returns a human label for the variant.
func (v Variant) String() string {
	switch v {
	case FourButton:
		return "4-button (digits 1-4)"
	case SixButton:
		return "6-button (digits 1-6)"
	case EightButton:
		return "8-button (digits 1-8)"
	default:
		return fmt.Sprintf("Variant(%d)", int(v))
	}
}

// fold maps any value into the 1..mod range the way the radio firmware does:
// 0 becomes 1, anything above mod has mod subtracted until it fits.
func fold(v, mod int) int {
	if v == 0 {
		return 1
	}
	return ((v-1)%mod)+1
}

// Code returns the 5-digit unlock code for the given last four serial digits.
// last4 must be exactly four decimal characters, e.g. "7885".
func Code(last4 string, v Variant) (string, error) {
	if v != FourButton && v != SixButton && v != EightButton {
		return "", fmt.Errorf("unknown variant %d", int(v))
	}
	if len(last4) != 4 {
		return "", fmt.Errorf("need exactly 4 digits, got %q", last4)
	}
	s := make([]int, 4)
	for i, c := range last4 {
		if c < '0' || c > '9' {
			return "", errors.New("serial digits must be 0-9")
		}
		s[i] = int(c - '0')
	}
	mod := int(v)

	// Treat the last four digits as a number, read its octal (3-bit) groups,
	// and re-pack them in a fixed order to get M.
	n := s[0]*1000 + s[1]*100 + s[2]*10 + s[3]
	b0 := n & 7
	b1 := (n >> 3) & 7
	b2 := (n >> 6) & 7
	b3 := (n >> 9) & 7
	b4 := (n >> 12) & 7
	m := (4096*b2 + 512*b4 + 64*b1 + 8*b0 + b3) % 10000

	// Five raw values: the digit sum, then M's four decimal digits.
	d := [5]int{s[0] + s[1] + s[2] + s[3], m / 1000, (m / 100) % 10, (m / 10) % 10, m % 10}
	r := [5]int{}
	for i, v := range d {
		r[i] = fold(v, mod)
	}
	// Anti-repeat rule: four identical trailing digits bumps the second one.
	if r[1] == r[2] && r[2] == r[3] && r[3] == r[4] {
		r[1] = fold(r[1]+1, mod)
	}

	out := make([]byte, 5)
	for i, v := range r {
		out[i] = byte('0' + v)
	}
	return string(out), nil
}
