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
			ctrlDInput(pages)
		}
		return event
	})
}

// ctrlDInput : Handler for Ctrl+D input.
func ctrlDInput(pages *tview.Pages) {
	// Create the modal to prompt user confirmation.
	var (
		buttonFn func()
		title    string
	)

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
