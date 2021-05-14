package pages

import (
	"fmt"
	"net/url"
	"strings"

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
		// Set up query parameters for the search.
		params := url.Values{}
		params.Add("limit", "75")
		setUpMangaListTable(pages, table, &params) // This will not set the titles, so we do it below.

		grid.SetTitle("Welcome to MangaDex, [red]Guest!")
		table.SetTitle("Recently updated manga")
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
	table.SetTitle("Your followed manga")

	// Colors for the table.
	titleColor := tcell.ColorLightGoldenrodYellow
	statusColor := tcell.ColorSaddleBrown

	// Set up table
	mangaTitleHeader := tview.NewTableCell("Manga"). // Manga header
								SetAlign(tview.AlignCenter).
								SetTextColor(titleColor).
								SetSelectable(false)
	statusHeader := tview.NewTableCell("Status"). // Status header
							SetAlign(tview.AlignCenter).
							SetTextColor(statusColor).
							SetSelectable(false)
	table.SetCell(0, 0, mangaTitleHeader). // Add the headers to the table
						SetCell(0, 1, statusHeader).
						SetFixed(1, 0) // This allows the table to show the header at all times.

	// GOROUTINE
	go func() {
		// Get user's followed manga.
		followedManga, err := g.Dex.GetUserFollowedMangaList(50, 0)
		if err != nil {
			// If error getting user's followed manga, we show a modal to the user indicating so.
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				ShowModal(pages, g.GenericAPIErrorModalID, "Error getting followed manga.", []string{"OK"},
					func(i int, label string) {
						pages.RemovePage(g.GenericAPIErrorModalID)
					})
			})
			return // We end immediately. No need to continue.
		}

		// When user presses ENTER on a manga row, they are redirected to the manga page.
		table.SetSelectedFunc(func(row, column int) {
			ShowMangaPage(pages, &(followedManga.Results[row-1]))
		})

		// Add each entry to the table.
		for i, mr := range followedManga.Results {
			// Create the manga title cell and fill it with the name of the manga.
			mtCell := tview.NewTableCell(fmt.Sprintf("%-50s", mr.Data.Attributes.Title["en"])).
				SetMaxWidth(50)
			mtCell.Color = titleColor

			// Create status cell and fill lit with the manga's current status.
			status := "-"
			if mr.Data.Attributes.Status != nil {
				status = strings.Title(*mr.Data.Attributes.Status)
			}
			sCell := tview.NewTableCell(fmt.Sprintf("%-15s", status)).
				SetMaxWidth(15)
			sCell.Color = statusColor

			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				table.SetCell(i+1, 0, mtCell).
					SetCell(i+1, 1, sCell)
			})
		}
	}()
}

// setUpMangaListTable : Set up table to show list of manga information.
func setUpMangaListTable(pages *tview.Pages, table *tview.Table, params *url.Values) {
	// For guest users, we fill the table with recently updated manga.
	titleColor := tcell.ColorOrange
	descColor := tcell.ColorLightGrey
	tagColor := tcell.ColorLightSteelBlue

	// Set up the table.
	mangaTitleHeader := tview.NewTableCell("Manga"). // Manga header
								SetAlign(tview.AlignCenter).
								SetTextColor(titleColor).
								SetSelectable(false)
	descHeader := tview.NewTableCell("Description"). // Description header
								SetAlign(tview.AlignCenter).
								SetTextColor(descColor).
								SetSelectable(false)
	tagHeader := tview.NewTableCell("Tags"). // Tag header
							SetAlign(tview.AlignCenter).
							SetTextColor(tagColor).
							SetSelectable(false)
	table.SetCell(0, 0, mangaTitleHeader). // Add headers to the table
						SetCell(0, 1, descHeader).
						SetCell(0, 2, tagHeader).
						SetFixed(1, 0) // This allows the table to show the header at all times.

	// GOROUTINE
	go func() {
		// Get manga list based on required query parameters.
		mangaList, err := g.Dex.MangaList(*params)
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

		// When user presses ENTER on a manga row, they are redirected to the manga page.
		table.SetSelectedFunc(func(row, column int) {
      if(len(mangaList.Results) != 0) {
		    ShowMangaPage(pages, &(mangaList.Results[row-1]))
      }
		})

		// Add each entry to the table.
		for i, mr := range mangaList.Results {
			// Create the manga title cell and fill it with the name of the manga.
			mtCell := tview.NewTableCell(fmt.Sprintf("%-40s", mr.Data.Attributes.Title["en"])).
				SetMaxWidth(40)
			mtCell.Color = titleColor

			// Create the description cell and fill it with the manga description.
			desc := tview.Escape(fmt.Sprintf("%-70s", mr.Data.Attributes.Description["en"]))
			descCell := tview.NewTableCell(desc).
				SetMaxWidth(70)
			descCell.Color = descColor

			// Create the tag cell and fill it with the manga tags.
			tags := make([]string, len(mr.Data.Attributes.Tags))
			for ti, tag := range mr.Data.Attributes.Tags {
				tags[ti] = tag.Attributes.Name["en"]
			}
			tagCell := tview.NewTableCell(strings.Join(tags, ", "))
			tagCell.Color = tagColor

			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				table.SetCell(i+1, 0, mtCell).
					SetCell(i+1, 1, descCell).
					SetCell(i+1, 2, tagCell)
			})
		}
	}()
}
