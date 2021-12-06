package ui

import (
	"github.com/darylhjd/mangadesk/app/core"
	"github.com/rivo/tview"
)

// SearchPage : This struct contains the search bar and the table of results
// for the search. This struct reuses the MainPage struct, specifically for the guest table.
type SearchPage struct {
	MainPage
	Form *tview.Form
}

// ShowSearchPage : Make the app show the search page.
func ShowSearchPage() {
	// Create the new search page
	searchPage := newSearchPage()

	core.App.TView.SetFocus(searchPage.Grid)
	core.App.PageHolder.AddAndSwitchToPage(SearchPageID, searchPage.Grid, true)
}

// newSearchPage : Creates a new SearchPage.
func newSearchPage() *SearchPage {
	var dimensions []int
	for i := 0; i < 15; i++ {
		dimensions = append(dimensions, -1)
	}
	grid := newGrid(dimensions, dimensions)
	// Set grid attributes
	grid.SetTitleColor(SearchPageGridTitleColor).
		SetBorderColor(SearchPageGridBorderColor).
		SetTitle("Search Manga. " +
			"[yellow]Press ↓ on search bar to switch to table. " +
			"[green]Press Tab on table to switch to search bar.").
		SetBorder(true)

	// Create table to show manga list.
	table := tview.NewTable()
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(SearchPageTableBorderColor).
		SetTitleColor(SearchPageTableTitleColor).
		SetTitle("The curious cat peeks into the database...🐈").
		SetBorder(true)

	// Create a form for the searching
	search := tview.NewForm()
	// Set form attributes
	search.SetButtonsAlign(tview.AlignLeft).
		SetLabelColor(SearchFormLabelColor)

	// Add search bar and result table to the grid. Search bar will have focus.
	grid.AddItem(search, 0, 0, 4, 15, 0, 0, false).
		AddItem(table, 4, 0, 11, 15, 0, 0, true)

	// Create the SearchPage.
	// We reuse the MainPage struct.
	searchPage := &SearchPage{
		MainPage: MainPage{
			Grid:  grid,
			Table: table,
		},
		Form: search,
	}

	// Add form fields
	search.AddInputField("Search Manga:", "", 0, nil, nil).
		AddCheckbox("Explicit Content?", false, nil).
		AddButton("Search", func() { // Search button.
			// Remove all current search results
			searchPage.Table.Clear()

			// When user presses button, we initiate the search.
			searchTerm := search.GetFormItemByLabel("Search Manga:").(*tview.InputField).GetText()
			exContent := search.GetFormItemByLabel("Explicit Content?").(*tview.Checkbox).IsChecked()
			go searchPage.setSearchTable(exContent, searchTerm)

			// Send focus to the search result table.
			core.App.TView.SetFocus(searchPage.Table)
		}).
		SetFocus(0) // Set focus to the title field.

	// Set handlers.
	searchPage.setHandlers()

	return searchPage
}

func (p *SearchPage) setSearchTable(exContent bool, searchTerm string) {
	p.MainPage.setGuestTable(true, exContent, searchTerm)
}
