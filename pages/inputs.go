package pages

/*
This file contains the input handlers for the pages.

The 2nd section of this page contains the logic for keybindings.
*/

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

// SetUniversalHandlers : Set input handlers for the app.
// List of input captures: Ctrl+L, Ctrl+K, Ctrl+S
func SetUniversalHandlers(pages *tview.Pages) {
	// Enable mouse.
	g.App.EnableMouse(true)

	// Set keyboard captures
	g.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlL: // Login/Logout
			ctrlLInput(pages)
		case tcell.KeyCtrlK: // Help page.
			ctrlKInput(pages)
		case tcell.KeyCtrlS: // Search page.
			ctrlSInput(pages)
		}
		return event
	})
}

// SetMainPageTableHandlers : Set input handler for main page table.
// List of input captures : Ctrl+F. Ctrl+B
func SetMainPageTableHandlers(cancel context.CancelFunc, pages *tview.Pages, mp *MainPage, searchTitle string, exContent bool) {
	// Pagination logic.
	mp.MangaListTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		changePage := false
		switch event.Key() {
		case tcell.KeyCtrlF: // User wants to go to next offset page.
			if mp.CurrentOffset+g.OffsetRange >= mp.MaxOffset {
				OKModal(pages, g.OffsetErrorModalID, "No more results to show.")
				break
			}
			mp.CurrentOffset += g.OffsetRange
			changePage = true
		case tcell.KeyCtrlB: // User wants to go back to previous offset page.
			if mp.CurrentOffset == 0 {
				OKModal(pages, g.OffsetErrorModalID, "Already on first page.")
				break
			}
			mp.CurrentOffset -= g.OffsetRange
			if mp.CurrentOffset < 0 {
				mp.CurrentOffset = 0
			}
			changePage = true
		}

		if changePage {
			cancel() // Cancel current goroutine.
			if mp.LoggedPage {
				mp.SetUpLoggedTable(pages)
			} else {
				// Get titles
				tableTitle := strings.SplitN(mp.MangaListTable.GetTitle(), ".", 2)[0] + "."
				mp.SetUpGenericTable(pages, tableTitle, searchTitle, exContent)
			}
		}
		return event
	})
}

// SetMangaPageHandlers : Set input handlers for the manga page.
// List of input captures : ESC
func SetMangaPageHandlers(cancel context.CancelFunc, pages *tview.Pages, grid *tview.Grid) {
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc: // User wants to go back.
			pages.RemovePage(g.MangaPageID)
			cancel()
		}
		return event
	})
}

// SetMangaPageTableHandlers : Set input handlers for the manga page table.
// List of input captures : Ctrl+E
func SetMangaPageTableHandlers(pages *tview.Pages, mangaPage *MangaPage, mr *mangodex.MangaResponse, chaps *[]mangodex.ChapterResponse) {
	// When user presses ENTER to confirm selected.
	mangaPage.ChapterTable.SetSelectedFunc(func(row, column int) {
		// We add the current selection if the there are no selected rows currently.
		if len(*mangaPage.Selected) == 0 {
			(*mangaPage.Selected)[row] = struct{}{}
		}
		// Show modal to confirm download.
		ShowModal(pages, g.DownloadChaptersModalID, "Download selection(s)?", []string{"Yes", "No"},
			func(i int, label string) {
				if label == "Yes" {
					// If user confirms to download, then we download the chapters.
					downloadChapters(pages, mangaPage, mr, chaps)
				}
				pages.RemovePage(g.DownloadChaptersModalID)
			})
	})

	// For custom input.
	mangaPage.ChapterTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlE: // User selects this manga row.
			ctrlEInput(mangaPage)
		case tcell.KeyCtrlA: // User wants to toggle select all.
			ctrlAInput(mangaPage, len(*chaps))
		}
		return event
	})
}

// SetSearchPageHandlers : Set input handlers for the search page.
// List of input captures : ESC, Tab, KeyDown
func SetSearchPageHandlers(pages *tview.Pages, searchPage *SearchPage) {
	// Set up input capture for the grid.
	searchPage.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc: // When user presses ESC, then we remove the Search page.
			pages.RemovePage(g.SearchPageID)
		case tcell.KeyTab: // When user presses Tab, they are sent back to the search form.
			g.App.SetFocus(searchPage.SearchForm)
		}
		return event
	})

	// Set up input capture for the search bar.
	searchPage.SearchForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown: // When user presses KeyDown, they are sent to the search results table.
			g.App.SetFocus(searchPage.MangaListTable)
		}
		return event
	})
}

// SetHelpPageHandlers : Set input handlers for the help page.
// List of input captures : ESC
func SetHelpPageHandlers(pages *tview.Pages, grid *tview.Grid) {
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc: // When user presses ESC, then we remove the Help Page.
			pages.RemovePage(g.HelpPageID)
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
			if err := g.Dex.Logout(); err != nil {
				OKModal(pages, g.LoginLogoutFailureModalID, "Error logging out!")
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

// ctrlKInput : Handler for Ctrl+K input.
// This shows the help page to the user.
func ctrlKInput(pages *tview.Pages) {
	ShowHelpPage(pages)
}

// ctrlSInput : Handler for Ctrl+S input.
// THis shows search page to the user.
func ctrlSInput(pages *tview.Pages) {
	// Do not allow when on login screen.
	if page, _ := pages.GetFrontPage(); page == g.LoginPageID {
		return
	}
	ShowSearchPage(pages)
}

// ctrlEInput : Handler for Ctrl+E input.
// This input enables user to select a chapter table row without activating the select action.
// This is done by using a int map to keep track of the selected row indexes.
func ctrlEInput(mangaPage *MangaPage) {
	// Get the current row (and col, but we do not need that)
	row, _ := mangaPage.ChapterTable.GetSelection()
	if _, ok := (*mangaPage.Selected)[row]; ok { // If the row already exists in the map, then we remove it!
		markChapterUnselected(mangaPage.ChapterTable, row)
		delete(*mangaPage.Selected, row)
	} else { // Else, we add the row into the map.
		markChapterSelected(mangaPage.ChapterTable, row)
		(*mangaPage.Selected)[row] = struct{}{}
	}
}

// ctrlAInput : Handler for Ctrl+A.
// This input enables users to select all chapters in the table.
func ctrlAInput(mangaPage *MangaPage, numChaps int) {
	// Note that the row is indexed with respect to the table.
	if mangaPage.SelectedAll { // If user has already selected all, then pressing again deselects all.
		mangaPage.Selected = &map[int]struct{}{} // Empty the map
		for r := 0; r < numChaps; r++ {
			markChapterUnselected(mangaPage.ChapterTable, r+1)
		}
	} else {
		for r := 0; r < numChaps; r++ {
			(*mangaPage.Selected)[r+1] = struct{}{}
			markChapterSelected(mangaPage.ChapterTable, r+1)
		}
	}
	mangaPage.SelectedAll = !mangaPage.SelectedAll
}
