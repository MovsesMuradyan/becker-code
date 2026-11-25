package becker

import "testing"

// Vectors confirmed against the original Becker generator: the eight BE2210
// forum pairs, three BE1650 pairs, and the user's own BE2210 serial.
func TestSixButton(t *testing.T) {
	cases := map[string]string{
		"8593": "15136", "6508": "11342", "4711": "14321", "2095": "41321",
		"9700": "43326", "6174": "61156", "6725": "24653", "2912": "21141",
		"6676": "11611", "5303": "53146", "4214": "55141", // BE1650
		"7885": "42311", // user's BE2210 (serial V1097885)
		"0000": "12111", // anti-repeat edge case
	}
	for in, want := range cases {
		got, err := Code(in, SixButton)
		if err != nil {
			t.Fatalf("Code(%q): %v", in, err)
		}
		if got != want {
			t.Errorf("Code(%q) = %s, want %s", in, got, want)
		}
	}
}

func TestOtherVariants(t *testing.T) {
	// Spot-checks captured from the reference generator.
	if got, _ := Code("0000", FourButton); got != "12111" {
		t.Errorf("4-button 0000 = %s, want 12111", got)
	}
	if got, _ := Code("0475", EightButton); got != "81888" {
		t.Errorf("8-button 0475 = %s, want 81888", got)
	}
}

func TestInputValidation(t *testing.T) {
	for _, bad := range []string{"", "123", "12345", "78a5", "-885"} {
		if _, err := Code(bad, SixButton); err == nil {
			t.Errorf("Code(%q) expected error, got nil", bad)
		}
	}
}
