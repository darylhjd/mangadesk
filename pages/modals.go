package pages

import (
	"github.com/darylhjd/mangadesk/core"
	"github.com/rivo/tview"
)

// OKModal : Creates a new modal with an "OK" acknowledgement button.
func OKModal(id, text string) *tview.Modal {
	modal := tview.NewModal()

	// Set modal attributes
	modal.SetText(text).
		AddButtons([]string{"OK"}).
		SetFocus(0).
		SetDoneFunc(func(_ int, _ string) {
			core.App.PageHolder.RemovePage(id)
		})

	return modal
}
