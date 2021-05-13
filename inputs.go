package main

import (
	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os"
)

// setUniversalInputCaptures : Set input handlers for the app.
// List of input captures: Ctrl+L
func setUniversalInputCaptures(pages *tview.Pages) {
	// Enable mouse.
	app.EnableMouse(true)

	// Set keyboard captures
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlL: // Login/Logout
			ctrlLInput(pages)
		case tcell.KeyCtrlH:
			ctrlHInput(pages)
		}
		return event
	})
}

// setMangaPageHandlers : Set input handlers for the manga page.
func setMangaPageHandlers(pages *tview.Pages, grid *tview.Grid) {
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			pages.RemovePage(MangaPageID)
		}
		return event
	})
}

// setMangaPageTableHandlers : Set input handlers for the manga page table.
func setMangaPageTableHandlers(pages *tview.Pages, table *tview.Table,
	selected *map[int]struct{}, mr *mangodex.MangaResponse, cl *mangodex.ChapterList) {
	// When user presses enter to confirm selected
	table.SetSelectedFunc(func(row, column int) {
		confirmChapterDownloads(pages, table, selected, row, mr, cl)
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlE:
			ctrlEInput(table, selected) // Check for updates for each manga.
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
		title = "Logout\nStored credentials will be deleted.\n\n"
		buttonFn = func() {
			// Attempt logout
			err := dex.Logout()
			if err != nil {
				ShowModal(pages, LoginLogoutFailureModalID, "Error logging out!", []string{"OK"},
					func(i int, label string) {
						pages.RemovePage(LoginLogoutFailureModalID)
					})
				return
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

// ctrlEInput : Handler for Ctrl+E input.
// This input enables user to select a chapter table row without activating the select action.
// This is done by using a int map to keep track of the selected row indexes.
func ctrlEInput(table *tview.Table, sRows *map[int]struct{}) {
	// Get the current row (and col, but we do not need that)
	row, _ := table.GetSelection()
	if _, ok := (*sRows)[row]; ok { // If the row already exists in the map, then we remove it!
		markChapterUnselected(table, row)
		delete(*sRows, row)
	} else { // Else, we add the row into the map.
		markChapterSelected(table, row)
		(*sRows)[row] = struct{}{}
	}
}

// ctrlHInput : Handler for Ctrl+H input.
// This shows the help page to the user.
func ctrlHInput(pages *tview.Pages) {
	ShowHelpPage(pages)
}
