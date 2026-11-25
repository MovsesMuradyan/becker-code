# Becker Radio Code

A small cross-platform desktop app (**macOS + Windows**) that recovers the
anti-theft unlock code for **Becker car radios** from the **last four digits of
the serial number**. Built in Go with a [Fyne](https://fyne.io) GUI.

Made for the common case of an owner who lost the code card after a battery
disconnect. Use it only on a radio you own.

---

## What it covers

The code format depends on how many preset buttons the radio uses to enter the
code. This tool implements the **4-, 6- and 8-button** families (all 5-digit).

| Variant | Code | Typical radios |
| --- | --- | --- |
| **6-button** (default) | 5 digits, each 1-6 | Mercedes Becker **"Special"**: BE2210, BE1650, BE1350 |
| 4-button | 5 digits, each 1-4 | 4-preset Becker units |
| 8-button | 5 digits, each 1-8 | 8-preset Becker units |

The **6-button** path is verified (see *Accuracy* below). The 4/8-button paths
use the same algorithm with a different modulus.

> The "Classic" line (e.g. BE2010) and newer navigation units (BE6xxx/BE7xxx,
> Traffic Pro, CDR) use different schemes and are **not** covered.

---

## Download & run

Prebuilt apps are produced by CI (GitHub Actions). From the repo's **Actions**
tab, open the latest **build** run and download the artifact for your OS, or
push a `v*` tag to attach them to a release.

- **macOS**: unzip and open `Becker Radio Code.app` (first launch: right-click →
  Open, since it is unsigned).
- **Windows**: run `Becker Radio Code.exe`.

Then type the last four serial digits, pick the radio type, and read the code.
Enter it on the radio's preset buttons and press **SC** to confirm.

---

## Build it yourself

Requires Go 1.23+. Fyne needs a C compiler and the usual GUI dev libraries
(see the [Fyne getting-started guide](https://docs.fyne.io/started/)).

```bash
go mod tidy
go run .                 # launch the GUI
go test ./becker/...     # run the verification tests

# package native apps with the Fyne tool:
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -os darwin  -name "Becker Radio Code" -icon Icon.png
fyne package -os windows -name "Becker Radio Code" -icon Icon.png
```

There is also a command-line version:

```bash
go run ./cmd/beckercli 7885        # 6-button (default)  -> 42311
go run ./cmd/beckercli 7885 4      # 4-button variant
```

---

## How the algorithm works

Using the last four serial digits `d0 d1 d2 d3` (example: `7885`):

1. Read them as a number `N` and split `N` into its five octal (3-bit) groups.
2. Re-pack those groups in a fixed order to get `M`
   (`M = 4096*b2 + 512*b4 + 64*b1 + 8*b0 + b3`).
3. Take `M mod 10000` and read its four decimal digits.
4. Prepend a fifth value equal to the **sum** of the four serial digits.
5. Fold each of the five values into the keypad range (1..N buttons): `0`
   becomes `1`, anything larger wraps down by the modulus.
6. If the last four output digits are identical, bump the second one (the
   firmware's anti-repeat rule).

Full code in [`becker/becker.go`](becker/becker.go).

## Accuracy

The implementation was checked against the original Becker generator for **all
10,000 possible inputs** on each variant (zero differences), and against real
serial-to-code pairs for **BE2210** (8/8) and **BE1650** (3/3).

## Disclaimer

For recovering the code of a radio you own. Not affiliated with Becker or
Mercedes-Benz. Provided as-is under the MIT License.
