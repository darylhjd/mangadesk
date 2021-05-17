package pages

import (
	"net/url"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

// ShowSearchPage : Show the search page to the user.
func ShowSearchPage(pages *tview.Pages) {
	// Create the base main grid.
	// 15x15 grid.
	var ga []int
	for i := 0; i < 15; i++ { // This is to create 15 grids.
		ga = append(ga, -1)
	}
	grid := tview.NewGrid().SetColumns(ga...).SetRows(ga...)
	// Set grid attributes
	grid.SetTitleColor(tcell.ColorOrange).
		SetBorderColor(tcell.ColorLightGrey).
		SetBorder(true)

	// Create table to show manga list.
	table := tview.NewTable()
	// Set table attributes
	table.SetSelectable(true, false). // Sets only the rows to be selectable
						SetSeparator('|').
						SetBordersColor(tcell.ColorGrey).
						SetTitleColor(tcell.ColorLightSkyBlue).
						SetBorder(true)

	// Create a form for the searching
	search := tview.NewForm()
	// Set form attributes
	search.SetButtonsAlign(tview.AlignLeft).
		SetLabelColor(tcell.ColorWhite).
		SetButtonBackgroundColor(tcell.ColorDodgerBlue)

	// Add form fields
	search.AddInputField("Search Manga:", "", 0, nil, nil). // Title field.
								AddButton("Search", func() { // Search button.
			// Remove all current search results
			table.Clear()

			// When user presses button, we initiate the search.
			searchTerm := search.GetFormItemByLabel("Search Manga:").(*tview.InputField).GetText()

			// Set up query parameters for the search.
			params := url.Values{}
			params.Add("limit", strconv.Itoa(g.OffsetRange))
			params.Add("title", searchTerm)
			title := "Search Results."
			setUpGenericMainPage(pages, grid, table, &params, title)

			// Set the correct titles, since the function sets the titles for the guest main page and not search.
			grid.SetTitle("Search Manga.")
			// table.SetTitle("[yellow]Press ↓ on search bar to switch to table. [green]Press Ctrl+Space on table to switch to search bar.")

			// Send focus to the search result table.
			g.App.SetFocus(table)
		}).SetFocus(0) // Set focus to the title field.

	// Set up input capture for the grid.
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc: // When user presses ESC, then we remove the Search page.
			pages.RemovePage(g.SearchPageID)
		case tcell.KeyCtrlSpace: // When user presses TAB, they are sent back to the search form.
			g.App.SetFocus(search)
		}
		return event
	})

	// Set up input capture for the search bar.
	search.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown:
			g.App.SetFocus(table)
		}
		return event
	})

	// Add search bar and result table to the grid. Search bar will have focus.
	grid.AddItem(search, 0, 0, 4, 15, 0, 0, false).
		AddItem(table, 4, 0, 11, 15, 0, 0, true)

	pages.AddPage(g.SearchPageID, grid, true, false)
	g.App.SetFocus(search)
	pages.SwitchToPage(g.SearchPageID)
}
