package core

import (
	"github.com/darylhjd/mangadesk/pages"
	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
)

// ShowLoginPage : Make the app show the login page.
func (m *MangaDesk) ShowLoginPage() {
	// Create the new login page
	loginPage := pages.NewLoginPage()

	m.ViewApp.SetFocus(loginPage.Form)
	m.PageHolder.AddAndSwitchToPage(pages.LoginPageID, loginPage.Grid, true)
}

// ShowMainPage : Make the app show the main page.
func (m *MangaDesk) ShowMainPage() {
	// Create the new main page
	mainPage := pages.NewMainPage()

	m.ViewApp.SetFocus(mainPage.Table)
	m.PageHolder.AddAndSwitchToPage(pages.MainPageID, mainPage.Grid, true)
}

func (m *MangaDesk) ShowMangaPage(manga *mangodex.Manga) {

}

// ShowModal : Make the app show a modal.
func (m *MangaDesk) ShowModal(id string, modal *tview.Modal) {
	m.PageHolder.AddPage(id, modal, true, true)
}
