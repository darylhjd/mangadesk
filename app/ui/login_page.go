package ui

import (
	"fmt"
	"github.com/darylhjd/mangadesk/app/core"
	"github.com/rivo/tview"
	"log"
)

// LoginPage : This struct contains the grid and form for the login page.
type LoginPage struct {
	Grid *tview.Grid
	Form *tview.Form
}

// ShowLoginPage : Make the app show the login page.
func ShowLoginPage() {
	// Create the new login page
	loginPage := newLoginPage()

	core.App.TView.SetFocus(loginPage.Grid)
	core.App.PageHolder.AddAndSwitchToPage(LoginPageID, loginPage.Grid, true)
}

// newLoginPage : Creates a new login page.
func newLoginPage() *LoginPage {
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
			loginPage.attemptLogin()
		}).
		AddButton("Guest", func() { // Guest button
			core.App.PageHolder.RemovePage(LoginPageID)
			ShowMainPage()
		})

	dimension := []int{0, 0, 0}
	grid := newGrid(dimension, dimension)

	grid.AddItem(form, 0, 0, 3, 3, 0, 0, true).
		AddItem(form, 1, 1, 1, 1, 32, 70, true)

	loginPage.Grid = grid
	loginPage.Form = form
	return loginPage
}

// attemptLogin : Attempts to log in with given form fields. If success, bring user to main page.
func (p *LoginPage) attemptLogin() {
	form := p.Form

	// Get username and password input.
	user := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	pwd := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	remember := form.GetFormItemByLabel("Remember Me").(*tview.Checkbox).IsChecked()

	// Attempt to log in to MangaDex API.
	if err := core.App.Client.Auth.Login(user, pwd); err != nil {
		log.Printf("Error trying to login: %ss\n", err.Error())
		modal := okModal(LoginLogoutFailureModalID, "Authentication failed.\nTry again!")
		ShowModal(LoginLogoutFailureModalID, modal)
		return
	}

	// Remember the user's login credentials if user wants it.
	if remember {
		if err := core.App.StoreCredentials(); err != nil {
			log.Println(fmt.Sprintf("Error storing credentials: %s", err.Error()))
			modal := okModal(StoreCredentialErrorModalID, "Failed to store login token.\nCheck logs for details.")
			ShowModal(StoreCredentialErrorModalID, modal)
		}
	}

	core.App.PageHolder.RemovePage(LoginPageID) // Remove the login page as we no longer need it.
	ShowMainPage()
}
