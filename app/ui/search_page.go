package ui

import (
	"context"
	"log"

	"github.com/darylhjd/mangadesk/app/core"
	"github.com/darylhjd/mangadesk/app/ui/utils"
	"github.com/rivo/tview"
)

// SearchPage : This struct contains the search bar and the table of results
// for the search. This struct reuses the MainPage struct, specifically for the guest table.
type SearchPage struct {
	MainPage
	Form *tview.Form
}

// SearchParams : Convenience struct to hold parameters for setting up a search table.
type SearchParams struct {
	term     string // The term to search for.
	explicit bool   // Whether to include explicit material in results.
}

// ShowSearchPage : Make the app show the search page.
func ShowSearchPage() {
	// Create the new search page
	searchPage := newSearchPage()

	core.App.PageHolder.AddAndSwitchToPage(utils.SearchPageID, searchPage.Grid, true)
	core.App.TView.SetFocus(searchPage.Form)
}

// newSearchPage : Creates a new SearchPage.
func newSearchPage() *SearchPage {
	var dimensions []int
	for i := 0; i < 15; i++ {
		dimensions = append(dimensions, -1)
	}
	grid := utils.NewGrid(dimensions, dimensions)
	// Set grid attributes
	grid.SetTitleColor(utils.SearchPageGridTitleColor).
		SetBorderColor(utils.SearchPageGridBorderColor).
		SetTitle("Search Manga. " +
			"[yellow]Press ↓ on search bar to switch to table. " +
			"[green]Press Tab on table to switch to search bar.").
		SetBorder(true)

	// Create table to show manga list.
	table := tview.NewTable()
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(utils.SearchPageTableBorderColor).
		SetTitleColor(utils.SearchPageTableTitleColor).
		SetTitle("The curious cat peeks into the database...🐈").
		SetBorder(true)

	// Create a form for the searching
	search := tview.NewForm()
	// Set form attributes
	search.SetButtonsAlign(tview.AlignLeft).
		SetLabelColor(utils.SearchFormLabelColor)

	// Add search bar and result table to the grid. Search bar will have focus.
	grid.AddItem(search, 0, 0, 4, 15, 0, 0, false).
		AddItem(table, 4, 0, 11, 15, 0, 0, true)

	// Create the SearchPage.
	// We reuse the MainPage struct.
	ctx, cancel := context.WithCancel(context.Background())
	searchPage := &SearchPage{
		MainPage: MainPage{
			Grid:  grid,
			Table: table,
			cWrap: &utils.ContextWrapper{
				Ctx:    ctx,
				Cancel: cancel,
			},
		},
		Form: search,
	}

	// Add form fields
	search.AddInputField("Search Manga:", "", 0, nil, nil).
		AddCheckbox("Explicit Content?", false, nil).
		AddButton("Search", func() { // Search button.
			// When user presses button, we initiate the search.
			searchTerm := search.GetFormItemByLabel("Search Manga:").(*tview.InputField).GetText()
			exContent := search.GetFormItemByLabel("Explicit Content?").(*tview.Checkbox).IsChecked()
			searchPage.setSearchTable(exContent, searchTerm)

			// Send focus to the search result table.
			core.App.TView.SetFocus(searchPage.Table)
		}).
		SetFocus(0) // Set focus to the title field.

	// Set handlers.
	searchPage.setHandlers()

	return searchPage
}

// setSearchTable : Sets the table for search results.
func (p *SearchPage) setSearchTable(exContent bool, searchTerm string) {
	log.Println("Setting new search results...")
	// Create the search param struct
	s := &SearchParams{
		term:     searchTerm,
		explicit: exContent,
	}
	go p.MainPage.setGuestTable(s)
}
