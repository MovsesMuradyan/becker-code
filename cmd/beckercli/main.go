// Command beckercli prints a Becker radio code for a serial's last four digits.
//
// Usage: beckercli <last4> [buttons]
//   beckercli 7885        // 6-button (default): 42311
//   beckercli 7885 4      // 4-button variant
package main

import (
	"fmt"
	"os"
	"strconv"

	"beckercode/becker"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: beckercli <last4-serial-digits> [buttons: 4|6|8]")
		os.Exit(2)
	}
	v := becker.SixButton
	if len(os.Args) > 2 {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, "buttons must be 4, 6 or 8")
			os.Exit(2)
		}
		v = becker.Variant(n)
	}
	code, err := becker.Code(os.Args[1], v)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println(code)
}
