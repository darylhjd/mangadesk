package ui

import (
	"github.com/darylhjd/mangadesk/app/core"
	"github.com/rivo/tview"
)

// ShowModal : Make the app show a modal.
func ShowModal(id string, modal *tview.Modal) {
	core.App.TView.SetFocus(modal)
	core.App.PageHolder.AddPage(id, modal, true, true)
}

// okModal : Creates a new modal with an "OK" acknowledgement button.
func okModal(id, text string) *tview.Modal {
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
