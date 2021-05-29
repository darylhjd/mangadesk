package pages

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

// TODO: Refactor this.

type SearchPage struct {
	MainPage
	SearchForm *tview.Form
}

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
	grid.SetTitleColor(g.SearchPageGridTitleColor).
		SetBorderColor(g.SearchPageGridBorderColor).
		SetTitle("Search Manga. " +
			"[yellow]Press ↓ on search bar to switch to table. " +
			"[green]Press Ctrl+Space on table to switch to search bar.").
		SetBorder(true)

	// Create table to show manga list.
	table := tview.NewTable()
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(g.SearchPageTableBorderColor).
		SetTitleColor(g.SearchPageTableTitleColor).
		SetTitle("The curious cat peeks into the database...🐈").
		SetBorder(true)

	// Create a form for the searching
	search := tview.NewForm()
	// Set form attributes
	search.SetButtonsAlign(tview.AlignLeft).
		SetLabelColor(g.SearchFormLabelColor)

	// Create the SearchPage.
	// We use the MainPage struct.
	searchPage := SearchPage{
		MainPage: MainPage{
			Grid:           grid,
			MangaListTable: table,
		},
		SearchForm: search,
	}

	// Add form fields
	search.AddInputField("Search Manga:", "", 0, nil, nil).
		AddButton("Search", func() { // Search button.
			// Remove all current search results
			searchPage.MainPage.MangaListTable.Clear()

			// When user presses button, we initiate the search.
			searchTerm := search.GetFormItemByLabel("Search Manga:").(*tview.InputField).GetText()
			searchPage.MainPage.SetUpGenericTable(pages, "Search Results.", 0, searchTerm)

			// Send focus to the search result table.
			g.App.SetFocus(searchPage.MainPage.MangaListTable)
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
