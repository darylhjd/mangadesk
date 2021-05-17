package pages

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

// ShowMainPage : Show the main page. Can be for logged user or guest user.
func ShowMainPage(pages *tview.Pages) {
	// Create the base main grid.
	// 15x15 grid.
	var ga []int
	for i := 0; i < 15; i++ { // This is to create 15 grids.
		ga = append(ga, -1)
	}
	grid := tview.NewGrid().SetColumns(ga...).SetRows(ga...)
	// Set grid attributes.
	grid.SetTitleColor(tcell.ColorOrange).
		SetBorderColor(tcell.ColorLightGrey).
		SetBorder(true)

	// Create the base main table.
	table := tview.NewTable()
	// Set table attributes
	table.SetSelectable(true, false). // Sets only the rows to be selectable
						SetSeparator('|').
						SetBordersColor(tcell.ColorGrey).
						SetTitleColor(tcell.ColorLightSkyBlue).
						SetBorder(true)

	// Add the table to the grid. Table spans the whole page.
	grid.AddItem(table, 0, 0, 15, 15, 0, 0, true)

	// Check if the user is logged in or not. Then, decide what to show for the main page.
	if g.Dex.IsLoggedIn() {
		setUpLoggedMainPage(pages, grid, table)
	} else {
		setUpGuestMainPage(pages, grid, table)
	}

	pages.AddPage(g.MainPageID, grid, true, false)
	g.App.SetFocus(grid)
	pages.SwitchToPage(g.MainPageID)
}

// setUpLoggedMainPage : Set up the main page for a logged user.
func setUpLoggedMainPage(pages *tview.Pages, grid *tview.Grid, table *tview.Table) {
	// For logged users, we fill the table with their followed manga.
	// Get user information
	username := "?"
	if u, err := g.Dex.GetLoggedUser(); err == nil {
		username = u.Data.Attributes.Username
	}
	grid.SetTitle(fmt.Sprintf("Welcome to MangaDex, [lightgreen]%s!", username))
	table.SetTitle("Your followed manga. Page 1")

	// Set up table
	mangaTitleHeader := tview.NewTableCell("Manga"). // Manga header
								SetAlign(tview.AlignCenter).
								SetTextColor(g.LoggedMainPageTitleColor).
								SetSelectable(false)
	statusHeader := tview.NewTableCell("Status"). // Status header
							SetAlign(tview.AlignCenter).
							SetTextColor(g.LoggedMainPageStatusColor).
							SetSelectable(false)
	table.SetCell(0, 0, mangaTitleHeader). // Add the headers to the table
						SetCell(0, 1, statusHeader).
						SetFixed(1, 0) // This allows the table to show the header at all times.

	go func() {
		setUpMangaListTable(pages, table, false, func() (*mangodex.MangaList, error) {
			return g.Dex.GetUserFollowedMangaList(100, 0)
		})
	}()
}

// setUpGuestMainPage : Set up a main page for a guest user.
func setUpGuestMainPage(pages *tview.Pages, grid *tview.Grid, table *tview.Table) {
	// For guest users, we fill the table with recently updated manga.
	grid.SetTitle("Welcome to MangaDex, [red]Guest!")
	table.SetTitle("Recently updated manga")

	// Set up the table.
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
	// Set up query parameters for the search.
	params := url.Values{}
	params.Add("limit", "50")
	go func() {
		setUpMangaListTable(pages, table, true, func() (*mangodex.MangaList, error) {
			return g.Dex.MangaList(params)
		})
	}()
}

// setUpMangaListTable : Set up the table for the manga list. Also used for the search page!
func setUpMangaListTable(pages *tview.Pages, table *tview.Table, full bool, f func() (*mangodex.MangaList, error)) {
	// Perform required search function for required manga list.
	mangaList, err := f()
	if err != nil {
		// If error getting the manga list, we show a modal to the user indicating so.
		g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
			ShowModal(pages, g.GenericAPIErrorModalID, "Error loading manga list.", []string{"OK"},
				func(i int, label string) {
					pages.RemovePage(g.GenericAPIErrorModalID)
				})
		})
		return // We end immediately. No need to continue.
	}

	// If no results, then tell user.
	if len(mangaList.Results) == 0 {
		noResCell := tview.NewTableCell("No results!").SetSelectable(false)
		g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
			table.SetCell(1, 0, noResCell)
		})
		return
	}

	// When user presses ENTER on a manga row, they are redirected to the manga page.
	table.SetSelectedFunc(func(row, column int) {
		// We do not need to worry about index out-of-range as we checked results earlier.
		ShowMangaPage(pages, &(mangaList.Results[row-1]))
	})

	// Add each entry to the table.
	if !full { // If not filling full information. This is usually for logged in main page.
		for i, mr := range mangaList.Results {
			// Create the manga title cell and fill it with the name of the manga.
			mtCell := tview.NewTableCell(fmt.Sprintf("%-50s", mr.Data.Attributes.Title["en"])).
				SetMaxWidth(50)
			mtCell.Color = g.LoggedMainPageTitleColor

			// Create status cell and fill lit with the manga's current status.
			status := "-"
			if mr.Data.Attributes.Status != nil {
				status = strings.Title(*mr.Data.Attributes.Status)
			}
			sCell := tview.NewTableCell(fmt.Sprintf("%-15s", status)).
				SetMaxWidth(15)
			sCell.Color = g.LoggedMainPageStatusColor

			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				table.SetCell(i+1, 0, mtCell).
					SetCell(i+1, 1, sCell)
			})
		}
	} else { // If showing full info, including tags and description.
		for i, mr := range mangaList.Results {
			// Create the manga title cell and fill it with the name of the manga.
			mtCell := tview.NewTableCell(fmt.Sprintf("%-40s", mr.Data.Attributes.Title["en"])).
				SetMaxWidth(40)
			mtCell.Color = g.GuestMainPageTitleColor

			// Create the description cell and fill it with the manga description.
			desc := tview.Escape(fmt.Sprintf("%-70s", mr.Data.Attributes.Description["en"]))
			descCell := tview.NewTableCell(desc).
				SetMaxWidth(70)
			descCell.Color = g.GuestMainPageDescColor

			// Create the tag cell and fill it with the manga tags.
			tags := make([]string, len(mr.Data.Attributes.Tags))
			for ti, tag := range mr.Data.Attributes.Tags {
				tags[ti] = tag.Attributes.Name["en"]
			}
			tagCell := tview.NewTableCell(strings.Join(tags, ", "))
			tagCell.Color = g.GuestMainPageTagColor

			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				table.SetCell(i+1, 0, mtCell).
					SetCell(i+1, 1, descCell).
					SetCell(i+1, 2, tagCell)
			})
		}
	}
}
