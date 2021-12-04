package pages

/*
Modal Page shows the modal on top of the current page (does not switch to, only shows on top of).
*/

import (
	"github.com/darylhjd/mangadesk/core"
	"github.com/rivo/tview"
)

// OKModal : Convenience function to show an acknowledgement modal.
func OKModal(m *core.MangaDesk, id, text string) *tview.Modal {
	modal := tview.NewModal()

	// Set modal attributes
	modal.SetText(text).
		AddButtons([]string{"OK"}).
		SetFocus(0).
		SetDoneFunc(func(_ int, _ string) {
			m.PageHolder.RemovePage(id)
		})

	return modal
}
