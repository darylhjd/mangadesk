package pages

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"

	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

// TODO: Refactor this.

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
	grid.SetTitleColor(g.MainPageGridTitleColor).
		SetBorderColor(g.MainPageGridBorderColor).
		SetBorder(true)

	// Create the base main table.
	table := tview.NewTable()
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(g.MainPageTableBorderColor).
		SetTitleColor(g.MainPageTableTitleColor).
		SetBorder(true)

	// Add the table to the grid. Table spans the whole page.
	grid.AddItem(table, 0, 0, 15, 15, 0, 0, true)

	// Check if the user is logged in or not. Then, decide what to show for the main page.
	if g.Dex.IsLoggedIn() {
		// If logged in, use the logged main page.
		offset := 0
		SetUpLoggedMainPage(pages, grid, table, &offset)
	} else {
		// If not logged in, use the generic main page.
		params := url.Values{}
		params.Add("limit", strconv.Itoa(g.OffsetRange))
		params.Add("offset", "0")
		title := "Recently updated manga."
		SetUpGenericMainPage(pages, grid, table, &params, title)
	}

	pages.AddPage(g.MainPageID, grid, true, false)
	g.App.SetFocus(grid)
	pages.SwitchToPage(g.MainPageID)
}

// SetUpLoggedMainPage : Set up the main page for a logged user.
func SetUpLoggedMainPage(pages *tview.Pages, grid *tview.Grid, table *tview.Table, offset *int) {
	// For logged users, we fill the table with their followed manga.
	// Get user information.
	username := "?"
	if u, err := g.Dex.GetLoggedUser(); err == nil {
		username = u.Data.Attributes.Username
	}
	welcome := "Welcome to MangaDex"
	if rand.Intn(100) < 3 { // 3% chance!
		welcome = "All according to keikaku (keikaku means plan)"
	}
	grid.SetTitle(fmt.Sprintf("%s, [lightgreen]%s!", welcome, username))

	// Set up table
	mangaTitleHeader := tview.NewTableCell("Title").
		SetAlign(tview.AlignCenter).
		SetTextColor(g.LoggedMainPageTitleColor).
		SetSelectable(false)
	pubStatusHeader := tview.NewTableCell("Pub. Status").
		SetAlign(tview.AlignLeft).
		SetTextColor(g.LoggedMainPagePubStatusColor).
		SetSelectable(false)
	table.SetCell(0, 0, mangaTitleHeader).
		SetCell(0, 1, pubStatusHeader).
		SetFixed(1, 0)

	go func() {
		// Perform required search function for required manga list.
		mangaList, err := g.Dex.GetUserFollowedMangaList(g.OffsetRange, *offset)
		if err != nil { // If error getting the manga list,
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				OKModal(pages, g.GenericAPIErrorModalID, "Error loading your followed manga.")
			})
			return // We end immediately. No need to continue.
		} else if len(mangaList.Results) == 0 { // If no followed manga.
			noResCell := tview.NewTableCell("You have no followed manga!").SetSelectable(false)
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				table.SetCell(1, 0, noResCell)
			})
			return // We end immediately. No need to continue.
		}

		// Set the title of the table.
		page := *offset/g.OffsetRange + 1
		firstEntry := *offset + 1
		lastEntry := page * g.OffsetRange
		if lastEntry > mangaList.Total {
			lastEntry = mangaList.Total
		}
		table.SetTitle(fmt.Sprintf("Your followed manga. Page %d (%d-%d).", page, firstEntry, lastEntry))

		// Set up input capture for the table. This allows for pagination logic.
		SetLoggedMainPageHandlers(pages, grid, table, mangaList, offset)

		// When user presses ENTER on a manga row, they are redirected to the manga page.
		table.SetSelectedFunc(func(row, column int) {
			// We do not need to worry about index out-of-range as we checked results earlier.
			ShowMangaPage(pages, &(mangaList.Results[row-1]))
		})

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
			sCell.Color = g.LoggedMainPagePubStatusColor

			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				table.SetCell(i+1, 0, mtCell).
					SetCell(i+1, 1, sCell)
			})
		}
	}()
}

// SetUpGenericMainPage : Set up a main page for a guest user.
func SetUpGenericMainPage(pages *tview.Pages, grid *tview.Grid, table *tview.Table, params *url.Values, title string) {
	// For guest users, we fill the table with recently updated manga.
	grid.SetTitle("Welcome to MangaDex, [red]Guest!")

	// Set up the table.
	mangaTitleHeader := tview.NewTableCell("Manga").
		SetAlign(tview.AlignCenter).
		SetTextColor(g.GuestMainPageTitleColor).
		SetSelectable(false)
	descHeader := tview.NewTableCell("Description").
		SetAlign(tview.AlignCenter).
		SetTextColor(g.GuestMainPageDescColor).
		SetSelectable(false)
	tagHeader := tview.NewTableCell("Tags").
		SetAlign(tview.AlignCenter).
		SetTextColor(g.GuestMainPageTagColor).
		SetSelectable(false)
	table.SetCell(0, 0, mangaTitleHeader).
		SetCell(0, 1, descHeader).
		SetCell(0, 2, tagHeader).
		SetFixed(1, 0)

	go func() { // Create the manga list table for a guest user.
		mangaList, err := g.Dex.MangaList(*params)
		if err != nil {
			// If error getting the manga list, we show a modal to the user indicating so.
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				OKModal(pages, g.GenericAPIErrorModalID, "Error loading manga list.")
			})
			return // We end immediately. No need to continue.
		}

		// Set the title of the table.
		offset, _ := strconv.Atoi(params.Get("offset"))
		page := offset/g.OffsetRange + 1
		firstEntry := offset + 1
		lastEntry := page * g.OffsetRange
		if lastEntry > mangaList.Total {
			lastEntry = mangaList.Total
		}
		table.SetTitle(fmt.Sprintf("%s Page %d (%d-%d).", title, page, firstEntry, lastEntry))

		// If no results, then tell user.
		if len(mangaList.Results) == 0 {
			noResCell := tview.NewTableCell("No results!").SetSelectable(false)
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				table.SetCell(1, 0, noResCell)
			})
			return // We end immediately. No need to continue.
		}

		// Set up input capture for the table. Allows for pagination logic.
		SetGuestMainPageHandlers(pages, grid, table, mangaList, params, title)

		// When user presses ENTER on a manga row, they are redirected to the manga page.
		table.SetSelectedFunc(func(row, column int) {
			// We do not need to worry about index out-of-range as we checked results earlier.
			ShowMangaPage(pages, &(mangaList.Results[row-1]))
		})

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
	}()
}
