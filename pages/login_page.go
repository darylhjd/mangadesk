package pages

/*
Login Page shows the form for the user to login or to continue as guest.
*/

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

// ShowLoginPage : Show login page.
func ShowLoginPage(pages *tview.Pages) {
	// Create the form.
	form := tview.NewForm()
	// Set form attributes.
	form.SetButtonsAlign(tview.AlignCenter).
		SetLabelColor(g.LoginFormLabelColor).
		SetTitle("Login to MangaDex").
		SetTitleColor(g.LoginPageTitleColor).
		SetBorder(true).
		SetBorderColor(g.LoginFormBorderColor)

	// Add form fields.
	form.AddInputField("Username", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddCheckbox("Remember Me", false, nil).
		AddButton("Login", func() {
			// Attempt login.
			if attemptLogin(pages, form) {
				// If we login successfully
				pages.RemovePage(g.LoginPageID) // Remove the login page as we no longer need it.
				ShowMainPage(pages)
			}
		}).
		AddButton("Guest", func() { // Guest button
			pages.RemovePage(g.LoginPageID)
			ShowMainPage(pages)
		})

	// Create a new grid for the form. This is to align the form to the centre.
	grid := tview.NewGrid().SetColumns(0, 0, 0).SetRows(0, 0, 0)
	grid.AddItem(form, 0, 0, 3, 3, 0, 0, true).
		AddItem(form, 1, 1, 1, 1, 32, 70, true)

	pages.AddPage(g.LoginPageID, grid, true, false)
	g.App.SetFocus(grid)
	pages.SwitchToPage(g.LoginPageID)
}

// attemptLogin : Attempts to login to MangaDex API with given form fields.
func attemptLogin(pages *tview.Pages, form *tview.Form) bool {
	// Get username and password input.
	u := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	p := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	remember := form.GetFormItemByLabel("Remember Me").(*tview.Checkbox).IsChecked()

	// Attempt to login to MangaDex API.
	if err := g.Dex.Login(u, p); err != nil {
		// If there was error during login, we create a Modal to tell the user that the login failed.
		OKModal(pages, g.LoginLogoutFailureModalID, "Authentication failed.\nTry again!")
		return false
	} else if remember && storeLoginDetails() { // Remember the user's login credentials if user wants it.
		// Error when storing credentials.
		OKModal(pages, g.StoreCredentialErrorModalID, "Error storing credentials.")
	}
	return true
}

// storeLoginDetails : Store the refresh token after logging in successfully if user wants to.
func storeLoginDetails() bool {
	// Store user credentials in `usr` folder. This is not (and should not be) changeable!
	if err := os.MkdirAll(g.ConfDir(), os.ModePerm); err != nil {
		return false
	}
	return ioutil.WriteFile(filepath.Join(g.ConfDir(), g.CredFileName), []byte(g.Dex.RefreshToken), os.ModePerm) != nil
}
