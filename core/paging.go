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

	m.ViewApp.SetFocus(loginPage.Grid)
	m.PageHolder.AddAndSwitchToPage(pages.LoginPageID, loginPage.Grid, true)
}

// ShowMainPage : Make the app show the main page.
func (m *MangaDesk) ShowMainPage() {
	// Create the new main page
	mainPage := pages.NewMainPage()

	m.ViewApp.SetFocus(mainPage.Grid)
	m.PageHolder.AddAndSwitchToPage(pages.MainPageID, mainPage.Grid, true)
}

// ShowMangaPage : Make the app show the manga page.
func (m *MangaDesk) ShowMangaPage(manga *mangodex.Manga) {
	mangaPage := pages.NewMangaPage(manga)

	m.ViewApp.SetFocus(mangaPage.Grid)
	m.PageHolder.AddAndSwitchToPage(pages.MangaPageID, mangaPage.Grid, true)
}

// ShowHelpPage : Make the app show the help page.
func (m *MangaDesk) ShowHelpPage() {
	helpPage := pages.NewHelpPage()

	m.ViewApp.SetFocus(helpPage.Grid)
	m.PageHolder.AddPage(pages.HelpPageID, helpPage.Grid, true, true)
}

// ShowModal : Make the app show a modal.
func (m *MangaDesk) ShowModal(id string, modal *tview.Modal) {
	m.ViewApp.SetFocus(modal)
	m.PageHolder.AddPage(id, modal, true, true)
}
