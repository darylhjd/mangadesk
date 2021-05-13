package pages

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ShowModal : Convenience function to create a modal.
func ShowModal(pages *tview.Pages, label, text string, buttons []string, f func(buttonIndex int, buttonLabel string)) {
	m := tview.NewModal()
	// Set modal attributes
	m.SetText(text).
		AddButtons(buttons).
		SetDoneFunc(f).
		SetFocus(0).
		SetBackgroundColor(tcell.ColorDarkSlateGrey)

	pages.AddPage(label, m, true, false)
	pages.ShowPage(label)
}
