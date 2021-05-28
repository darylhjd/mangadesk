package pages

import (
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

// ShowModal : Create a modal and show it with required buttons and done function.
func ShowModal(pages *tview.Pages, modalID, text string, buttons []string, f func(buttonIndex int, buttonLabel string)) {
	m := tview.NewModal()
	// Set modal attributes
	m.SetText(text).
		AddButtons(buttons).
		SetDoneFunc(f).
		SetFocus(0).
		SetBackgroundColor(g.ModalColor)

	pages.AddPage(modalID, m, true, false)
	pages.ShowPage(modalID)
}

// OKModal : Convenience function to show an acknowledgement modal.
func OKModal(pages *tview.Pages, modalID, text string) {
	ShowModal(pages, modalID, text, []string{"OK"}, func(buttonIndex int, buttonLabel string) {
		pages.RemovePage(modalID)
	})
}
