package pages

/*
Search Page shows the interface for searching the MangaDex database.
 */

import (
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

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
			CurrentOffset:  0,
			MaxOffset:      0,
		},
		SearchForm: search,
	}

	// Add form fields
	search.AddInputField("Search Manga:", "", 0, nil, nil).
		AddButton("Search", func() { // Search button.
			// Remove all current search results
			searchPage.MangaListTable.Clear()

			// When user presses button, we initiate the search.
			searchTerm := search.GetFormItemByLabel("Search Manga:").(*tview.InputField).GetText()
			searchPage.MainPage.SetUpGenericTable(pages, "Search Results.", searchTerm)

			// Send focus to the search result table.
			g.App.SetFocus(searchPage.MangaListTable)
		}).SetFocus(0) // Set focus to the title field.

	// Set up input capture for the search page.
	SetSearchPageHandlers(pages, &searchPage)

	// Add search bar and result table to the grid. Search bar will have focus.
	grid.AddItem(search, 0, 0, 4, 15, 0, 0, false).
		AddItem(table, 4, 0, 11, 15, 0, 0, true)

	pages.AddPage(g.SearchPageID, grid, true, false)
	g.App.SetFocus(search)
	pages.SwitchToPage(g.SearchPageID)
}
