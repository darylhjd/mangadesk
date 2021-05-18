package pages

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

// ShowLoginPage : Show login page.
func ShowLoginPage(pages *tview.Pages) {
	// Create the form
	form := tview.NewForm()
	// Set form attributes.
	form.SetButtonsAlign(tview.AlignCenter).
		SetLabelColor(tcell.ColorWhite).
		SetTitle("Login to MangaDex").
		SetTitleColor(tcell.ColorOrange).
		SetBorder(true).
		SetBorderColor(tcell.ColorGrey)
	// Add form fields.
	form.AddInputField("Username", "", 0, nil, nil). // Username field
								AddPasswordField("Password", "", 0, '*', nil). // Password field
								AddCheckbox("Remember Me", false, nil).        // Remember Me field.
								AddButton("Login", func() {                    // Login button
			// Attempt login.
			if attemptLogin(pages, form) {
				// If we login successfully
				pages.RemovePage(g.LoginPageID) // Remove the login page as we no longer need it.
				ShowMainPage(pages)             // Switch to main page.
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
		ShowModal(pages, g.LoginLogoutFailureModalID, "Authentication failed\nTry again!", []string{"OK"},
			func(i int, l string) {
				pages.RemovePage(g.LoginLogoutFailureModalID) // Remove the modal once user acknowledge.
			})
		return false
	}
	// If successful login.
	// Remember the user's login credentials if user wants it.
	if remember && storeLoginDetails() {
		// Error when storing credentials.
		ShowModal(pages, g.StoreCredentialErrorModalID, "Error storing credentials.", []string{"OK"},
			func(i int, l string) {
				pages.RemovePage(g.StoreCredentialErrorModalID)
			})
	}
	return true
}

// storeLoginDetails : Store the refresh token after logging in successfully if user wants to.
func storeLoginDetails() bool {
	// Store user credentials in `usr` folder. This is not (and should not be) changeable!
	if err := os.MkdirAll(g.UsrDir, os.ModePerm); err != nil {
		return false
	}

	// Create file to store the token.
	return ioutil.WriteFile(filepath.Join(g.UsrDir, g.CredFileName), []byte(g.Dex.RefreshToken), os.ModePerm) != nil
}
