package core

import (
	"github.com/darylhjd/mangadesk/pages"
	"github.com/rivo/tview"
)

// ShowLoginPage : Make the app show the login page.
func (m *MangaDesk) ShowLoginPage() {
	// Create the new login page
	loginPage := pages.NewLoginPage(m)

	m.ViewApp.SetFocus(loginPage.Grid)
	m.PageHolder.AddAndSwitchToPage(pages.LoginPageID, loginPage.Grid, true)
}

func (m *MangaDesk) ShowMainPage() {

}

func (m *MangaDesk) ShowModal(id string, modal *tview.Modal) {
	m.PageHolder.AddPage(id, modal, true, true)
}
