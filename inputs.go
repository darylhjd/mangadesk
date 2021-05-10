package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
)

// SetInputCaptures : Set input handlers for the app.
func SetInputCaptures(pages *tview.Pages) {
	// Enable mouse.
	app.EnableMouse(true)

	// Set keyboard captures
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlL: // Login/Logout
			ctrlLInput(pages)
		}
		return event
	})
}

// ctrlLInput : Handler for Ctrl+D input.
func ctrlLInput(pages *tview.Pages) {
	// Do not allow pop up when on login screen.
	if page, _ := pages.GetFrontPage(); page == LoginPageID {
		return
	}

	// Create the modal to prompt user confirmation.
	var (
		buttonFn func()
		title    string
	)

	// This will decide the function of the modal.
	switch dex.IsLoggedIn() {
	case true: // User wants to logout.
		title = "Logout\n"
		buttonFn = func() {
			// Attempt logout
			if err := dex.Logout(); err != nil {
				panic(err)
			}
			// Remove the credentials file
			_ = os.Remove(credFile)
			// Then we redirect the user to the guest main page
			ShowMainPage(pages)
		}
	case false: // User wants to login.
		title = "Login\n"
		buttonFn = func() {
			ShowLoginPage(pages)
		}
	}

	// Create the modal for the user.
	lModal := CreateModal(title+"Are you sure?", []string{"Yes", "No"}, func(i int, label string) {
		// If user confirms the modal.
		if label == "Yes" {
			buttonFn()
		}
		// We remove the modal from the page.
		pages.RemovePage(LoginLogoutCfmModalID)
	})

	// Add and show the modal on top of the current screen.
	pages.AddPage(LoginLogoutCfmModalID, lModal, true, false)
	pages.ShowPage(LoginLogoutCfmModalID)
}
