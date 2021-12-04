package pages

import (
	"fmt"
	"github.com/rivo/tview"
	"log"

	"github.com/darylhjd/mangadesk/core"
)

// LoginPage : This struct contains the grid and form for the login page.
type LoginPage struct {
	Grid *tview.Grid
	Form *tview.Form
}

// NewLoginPage : Creates a new login page.
func NewLoginPage(m *core.MangaDesk) *LoginPage {
	// Create the LoginPage
	loginPage := &LoginPage{}

	form := tview.NewForm()

	// Set form attributes.
	form.SetButtonsAlign(tview.AlignCenter).
		SetLabelColor(LoginFormLabelColor).
		SetTitle("Login to MangaDex").
		SetTitleColor(LoginPageTitleColor).
		SetBorder(true).
		SetBorderColor(LoginFormBorderColor)

	// Add form fields.
	form.AddInputField("Username", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddCheckbox("Remember Me", false, nil).
		AddButton("Login", func() {
			loginPage.attemptLogin(m)
		}).
		AddButton("Guest", func() { // Guest button
			m.PageHolder.RemovePage(LoginPageID)
			m.ShowMainPage()
		})

	dimension := []int{0, 0, 0}
	grid := NewGrid(dimension, dimension)

	grid.AddItem(form, 0, 0, 3, 3, 0, 0, true).
		AddItem(form, 1, 1, 1, 1, 32, 70, true)

	loginPage.Grid = grid
	loginPage.Form = form
	return loginPage
}

// attemptLogin : Attempts to log in with given form fields. If success, bring user to main page.
func (f *LoginPage) attemptLogin(m *core.MangaDesk) {
	form := f.Form

	// Get username and password input.
	u := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	p := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	remember := form.GetFormItemByLabel("Remember Me").(*tview.Checkbox).IsChecked()

	// Attempt to log in to MangaDex API.
	if err := m.Client.Auth.Login(u, p); err != nil {
		modal := OKModal(m, LoginLogoutFailureModalID, "Authentication failed.\nTry again!")
		m.ShowModal(LoginLogoutFailureModalID, modal)
		return
	}

	// Remember the user's login credentials if user wants it.
	if remember {
		if err := m.StoreCredentials(); err != nil {
			log.Println(fmt.Sprintf("Error storing credentials: %s\n", err.Error()))
		}
	}

	m.PageHolder.RemovePage(LoginPageID) // Remove the login page as we no longer need it.
	m.ShowMainPage()
}
