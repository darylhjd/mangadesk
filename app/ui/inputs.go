package ui

import (
	"context"
	"log"
	"math"

	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/darylhjd/mangadesk/app/core"
)

// SetUniversalHandlers : Set universal inputs for the app.
func SetUniversalHandlers() {
	// Enable mouse inputs.
	core.App.TView.EnableMouse(true)

	// Set universal keybindings
	core.App.TView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlL: // Login/Logout
			ctrlLInput()
		case tcell.KeyCtrlK: // Help page.
			ctrlKInput()
		case tcell.KeyCtrlS: // Search page.
			ctrlSInput()
		case tcell.KeyCtrlC: // Ctrl-C interrupt.
			ctrlCInput()
		}
		return event // Forward the event to the actual current primitive.
	})
}

// ctrlLInput : Enables user to toggle login/logout.
func ctrlLInput() {
	log.Println("Received toggle login/logout event.")
	// Do not allow pop up when on login screen.
	if page, _ := core.App.PageHolder.GetFrontPage(); page == LoginPageID {
		return
	}

	// Create the modal to prompt user confirmation.
	var modal *tview.Modal
	// Decide whether the modal is to log in or logout.
	switch core.App.Client.Auth.IsLoggedIn() {
	case true:
		text := "Logout?\nStored credentials will be deleted."
		modal = confirmModal(LoginLogoutCfmModalID, text, "Logout", func() {
			// Attempt to logout
			if err := core.App.Client.Auth.Logout(); err != nil {
				okM := okModal(LoginLogoutFailureModalID, "Error logging out!")
				ShowModal(LoginLogoutFailureModalID, okM)
				return
			}
			// If logged out successfully, then delete stored credentials and direct user to main page (guest).
			core.App.DeleteCredentials()
			ShowMainPage()
		})
	case false:
		text := "Login?"
		modal = confirmModal(LoginLogoutCfmModalID, text, "Login", func() {
			ShowLoginPage()
		})
	}

	ShowModal(LoginLogoutCfmModalID, modal)
}

// ctrlKInput : Shows the help page to the user.
func ctrlKInput() {
	ShowHelpPage()
}

// ctrlSInput : Shows search page to the user.
func ctrlSInput() {
	// Do not allow when on login screen.
	if page, _ := core.App.PageHolder.GetFrontPage(); page == LoginPageID {
		return
	}
	ShowSearchPage()
}

// ctrlCInput : Sends an interrupt signal to the application to stop.
func ctrlCInput() {
	log.Println("TView stopped by Ctrl-C interrupt.")
	core.App.TView.Stop()
}

// setHandlers : Set handlers for the help page.
func (p *HelpPage) setHandlers() {
	p.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			core.App.PageHolder.RemovePage(HelpPageID)
		}
		return event
	})
}

// setHandlers : Set handlers for the search page.
func (p *SearchPage) setHandlers() {
	p.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc: // When user presses ESC, then we remove the Search page.
			core.App.PageHolder.RemovePage(SearchPageID)
		case tcell.KeyTab: // When user presses Tab, they are sent back to the search form.
			core.App.TView.SetFocus(p.Form)
		}
		return event
	})

	// Set up input capture for the search bar.
	p.Form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown: // When user presses KeyDown, they are sent to the search results table.
			core.App.TView.SetFocus(p.Table)
		}
		return event
	})
}

// setHandlers : Set handlers for the main page.
func (p *MainPage) setHandlers(cancel context.CancelFunc, searchParams *SearchParams) {
	p.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		var reload bool
		switch event.Key() {
		// User wants to go to the next offset page.
		case tcell.KeyCtrlF:
			if p.CurrentOffset+offsetRange >= p.MaxOffset {
				modal := okModal(OffsetErrorModalID, "No more results to show.")
				ShowModal(OffsetErrorModalID, modal)
				return event
			}
			// Update the new offset
			reload = true
			p.CurrentOffset += offsetRange
		case tcell.KeyCtrlB:
			if p.CurrentOffset == 0 {
				modal := okModal(OffsetErrorModalID, "Already on first page.")
				ShowModal(OffsetErrorModalID, modal)
				return event
			}
			reload = true
			// Update the new offset
			p.CurrentOffset = int(math.Max(0, float64(p.CurrentOffset-offsetRange)))
		}

		if reload {
			// Cancel any current loading, and create a new one.
			cancel()
			if searchParams != nil {
				go p.setGuestTable(searchParams)
			} else if !core.App.Client.Auth.IsLoggedIn() {
				go p.setGuestTable(nil)
			} else {
				go p.setLoggedTable()
			}
		}
		return event
	})

	p.Table.SetSelectedFunc(func(row, _ int) {
		log.Printf("Selected row %d on main page.\n", row)
		mangaRef := p.Table.GetCell(row, 0).GetReference()
		if mangaRef == nil {
			return
		} else if manga, ok := mangaRef.(*mangodex.Manga); ok {
			ShowMangaPage(manga)
		}
	})
}

// setHandlers : Set handlers for the manga page.
func (p *MangaPage) setHandlers(cancel context.CancelFunc) {
	p.Grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			cancel()
			core.App.PageHolder.RemovePage(MangaPageID)
		}
		return event
	})

	p.Table.SetSelectedFunc(func(row, _ int) {
		// We add the current selection if the there are no selected rows currently.
		if !p.sWrap.hasSelections() {
			p.sWrap.addSelection(row)
		}
		log.Println("Creating and showing confirm download modal...")
		modal := confirmModal(DownloadChaptersModalID, "Download chapter(s)?", "Yes", func() {
			// Create a copy of the selection.
			selected := p.sWrap.copySelection()
			// Set selected rows as unselected.
			for row := range selected {
				p.markUnselected(row)
			}
			// Download selected chapters.
			go p.downloadChapters(selected, 0)
		})
		ShowModal(DownloadChaptersModalID, modal)
	})

	p.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlE: // User selects this manga row.
			p.ctrlEInput()
		case tcell.KeyCtrlA: // User wants to toggle select all.
			p.ctrlAInput()
		}
		return event
	})
}

// ctrlEInput : Enables user to select a chapter table row without activating the select action.
func (p *MangaPage) ctrlEInput() {
	row, _ := p.Table.GetSelection()
	// If the row is already in the selection, we deselect. Else, we add.
	if p.sWrap.hasSelection(row) {
		p.markUnselected(row)
	} else {
		p.markSelected(row)
	}
}

// ctrlAInput : Enables user to select/deselect ALL chapters at once.
func (p *MangaPage) ctrlAInput() {
	// Toggle selection.
	p.markAll(!p.sWrap.all)
	p.sWrap.all = !p.sWrap.all
}
