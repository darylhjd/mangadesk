package pages

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

// SetUniversalInputCaptures : Set input handlers for the app.
// List of input captures: Ctrl+L
func SetUniversalInputCaptures(pages *tview.Pages) {
	// Enable mouse.
	g.App.EnableMouse(true)

	// Set keyboard captures
	g.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlL: // Login/Logout
			ctrlLInput(pages)
		case tcell.KeyCtrlH: // Help page.
			ctrlHInput(pages)
		case tcell.KeyCtrlS: // Search page.
			ctrlSInput(pages)
		}
		return event
	})
}

// SetMangaPageHandlers : Set input handlers for the manga page.
// List of input captures : ESC
func SetMangaPageHandlers(pages *tview.Pages, grid *tview.Grid) {
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc: // User wants to go back.
			pages.RemovePage(g.MangaPageID)
		}
		return event
	})
}

// SetMangaPageTableHandlers : Set input handlers for the manga page table.
// List of input captures : Ctrl+E
func SetMangaPageTableHandlers(table *tview.Table, selected *map[int]struct{}) {
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlE: // User selects this manga row.
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
	if page, _ := pages.GetFrontPage(); page == g.LoginPageID {
		return
	}

	// Create the modal to prompt user confirmation.
	var (
		buttonFn func()
		title    string
	)

	// This will decide the function of the modal by checking whether user is logged in or out.
	switch g.Dex.IsLoggedIn() {
	case true: // User wants to logout.
		title = "Logout\nStored credentials will be deleted.\n\n"
		buttonFn = func() { // Set the function.
			// Attempt logout
			err := g.Dex.Logout()
			if err != nil {
				ShowModal(pages, g.LoginLogoutFailureModalID, "Error logging out!", []string{"OK"},
					func(i int, label string) {
						pages.RemovePage(g.LoginLogoutFailureModalID)
					})
				return
			}
			// Remove the credentials file, but we ignore errors.
			_ = os.Remove(filepath.Join(g.UsrDir, g.CredFileName))
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
	ShowModal(pages, g.LoginLogoutCfmModalID, title+"Are you sure?", []string{"Yes", "No"},
		func(i int, label string) {
			// If user confirms the modal.
			if label == "Yes" {
				buttonFn()
			}
			// We remove the modal from the page.
			pages.RemovePage(g.LoginLogoutCfmModalID)
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

// ctrlSInput : Handler for Ctrl+S input.
// THis shows search page to the user.
func ctrlSInput(pages *tview.Pages) {
	if page, _ := pages.GetFrontPage(); page == g.LoginPageID {
		return
	}
	ShowSearchPage(pages)
}
