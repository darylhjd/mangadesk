package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
)

// setUniversalInputCaptures : Set input handlers for the app.
func setUniversalInputCaptures(pages *tview.Pages) {
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

// setMainPageTableInputCaptures : Set input handlers for the main page table.
func setMainPageTableInputCaptures(table *tview.Table, sRows *map[int]struct{}) {
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlE:
			ctrlEInput(table, sRows)
		}
		return event
	})
}

/*
It should be good design (based on what I know anyway) to not have overlapping input handlers for
different screens. As such, we try to only have one action for each input.
*/

// ctrlLInput : Handler for Ctrl+L input.
// This input enables user to toggle login/logout.
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
	ShowModal(pages, LoginLogoutCfmModalID, title+"Are you sure?", []string{"Yes", "No"},
		func(i int, label string) {
			// If user confirms the modal.
			if label == "Yes" {
				buttonFn()
			}
			// We remove the modal from the page.
			pages.RemovePage(LoginLogoutCfmModalID)
		})
}

// ctrlEInput() : Handler for Ctrl+E input.
// This input enables user to select a main page table row without activating the select action.
// This is done by using a int array to keep track of the selected row indexes.
func ctrlEInput(table *tview.Table, sRows *map[int]struct{}) {
	// Get the current row (and col, but we do not need that)
	row, _ := table.GetSelection()
	// If the row already exists in the map, then we remove it!
	if _, ok := (*sRows)[row]; ok {
		markChapterUnselected(table, row, tcell.ColorBlack, tcell.ColorWhite)
		delete(*sRows, row)
	} else { // Else, we add the row into the map.
		markChapterSelected(table, row, tcell.ColorLightSkyBlue, tcell.ColorBlack)
		(*sRows)[row] = struct{}{}
	}
}
