package pages

import (
	"github.com/darylhjd/mangodex"
	"net/url"

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
		SetTitle("Search for Manga").
		SetBorderColor(tcell.ColorLightGrey).
		SetBorder(true)

	// Set input handlers for this case
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc: // When user presses ESC, then we remove the Search page.
			pages.RemovePage(g.SearchPageID)
		}
		return event
	})

	// Create table to show manga list.
	table := tview.NewTable()
	// Set table attributes
	table.SetSelectable(true, false). // Sets only the rows to be selectable
						SetSeparator('|').
						SetBordersColor(tcell.ColorGrey).
						SetTitleColor(tcell.ColorLightSkyBlue).
						SetTitle("[yellow]Press ↓ on search bar to switch to table. [green]Press Tab on table to switch to search bar.").
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
			mangaTitleHeader := tview.NewTableCell("Manga"). // Manga header
										SetAlign(tview.AlignCenter).
										SetTextColor(g.GuestMainPageTitleColor).
										SetSelectable(false)
			descHeader := tview.NewTableCell("Description"). // Description header
										SetAlign(tview.AlignCenter).
										SetTextColor(g.GuestMainPageDescColor).
										SetSelectable(false)
			tagHeader := tview.NewTableCell("Tags"). // Tag header
									SetAlign(tview.AlignCenter).
									SetTextColor(g.GuestMainPageTagColor).
									SetSelectable(false)
			table.SetCell(0, 0, mangaTitleHeader). // Add headers to the table
								SetCell(0, 1, descHeader).
								SetCell(0, 2, tagHeader).
								SetFixed(1, 0) // This allows the table to show the header at all times.

			// When user presses button, we initiate the search.
			searchTerm := search.GetFormItemByLabel("Search Manga:").(*tview.InputField).GetText()

			// Set up query parameters for the search.
			params := url.Values{}
			params.Add("limit", "100")
			params.Add("title", searchTerm)
			go func() {
				setUpMangaListTable(pages, table, true, func() (*mangodex.MangaList, error) {
					return g.Dex.MangaList(params)
				})
			}()
			// Send focus to the search result table.
			g.App.SetFocus(table)
		}).SetFocus(0) // Set focus to the title field.

	// Set up input capture for the table.
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab: // When user presses TAB, they are sent back to the search form.
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
