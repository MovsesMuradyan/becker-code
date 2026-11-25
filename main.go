// Command becker-code is a small cross-platform desktop app (macOS and Windows)
// that turns the last four digits of a Becker radio serial into its unlock code.
package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"beckercode/becker"
)

func main() {
	a := app.New()
	w := a.NewWindow("Becker Radio Code")
	w.Resize(fyne.NewSize(420, 360))

	title := widget.NewLabelWithStyle("Becker Radio Code", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	intro := widget.NewLabel("Enter the last 4 digits of the radio serial number.")
	intro.Wrapping = fyne.TextWrapWord
	intro.Alignment = fyne.TextAlignCenter

	serial := widget.NewEntry()
	serial.SetPlaceHolder("e.g. 7885")

	variants := map[string]becker.Variant{
		"6-button / Mercedes Special (BE2210, BE1650)": becker.SixButton,
		"4-button":                                     becker.FourButton,
		"8-button":                                     becker.EightButton,
	}
	variantSel := widget.NewSelect([]string{
		"6-button / Mercedes Special (BE2210, BE1650)",
		"4-button",
		"8-button",
	}, nil)
	variantSel.SetSelectedIndex(0)

	result := canvas.NewText("-----", theme.ForegroundColor())
	result.TextSize = 46
	result.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	result.Alignment = fyne.TextAlignCenter

	status := canvas.NewText("", color.NRGBA{R: 0xC0, G: 0x30, B: 0x30, A: 0xFF})
	status.Alignment = fyne.TextAlignCenter

	recompute := func() {
		v := variants[variantSel.Selected]
		code, err := becker.Code(serial.Text, v)
		if err != nil {
			result.Text = "-----"
			if serial.Text != "" {
				status.Text = "Enter exactly 4 digits (0-9)."
			} else {
				status.Text = ""
			}
			result.Refresh()
			status.Refresh()
			return
		}
		result.Text = code
		status.Text = "Enter on preset buttons, then press SC."
		result.Refresh()
		status.Refresh()
	}

	serial.OnChanged = func(string) { recompute() }
	variantSel.OnChanged = func(string) { recompute() }

	copyBtn := widget.NewButton("Copy code", func() {
		if result.Text != "-----" {
			w.Clipboard().SetContent(result.Text)
		}
	})

	form := container.NewVBox(
		title,
		intro,
		widget.NewForm(
			widget.NewFormItem("Last 4 digits", serial),
			widget.NewFormItem("Radio type", variantSel),
		),
		widget.NewSeparator(),
		result,
		status,
		copyBtn,
	)

	w.SetContent(container.NewPadded(form))
	recompute()
	w.ShowAndRun()
}
