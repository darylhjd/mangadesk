package pages

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"

	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

type MainPage struct {
	Grid           *tview.Grid  // The page grid.
	MangaListTable *tview.Table // The table that contains the list of manga.
}

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

	// Create the MainPage.
	mainPage := MainPage{
		Grid:           grid,
		MangaListTable: table,
	}

	// Decide what to show for the main page.
	if g.Dex.IsLoggedIn() {
		mainPage.SetUpLoggedPage(pages)
	} else {
		mainPage.SetUpGenericPage(pages, "Welcome to MangaDex, [red]Guest!", "Popular Manga.")
	}

	pages.AddPage(g.MainPageID, grid, true, false)
	g.App.SetFocus(grid)
	pages.SwitchToPage(g.MainPageID)
}

// SetUpLoggedPage : Readies the MainPage to be a logged page.
// This will also call SetUpLoggedTable.
func (mp *MainPage) SetUpLoggedPage(pages *tview.Pages) {
	// Get user information
	username := "?"
	if u, err := g.Dex.GetLoggedUser(); err == nil {
		username = u.Data.Attributes.Username
	}
	welcome := "Welcome to MangaDex"
	if rand.Intn(100) < 3 { // 3% chance!
		welcome = "All according to keikaku (keikaku means plan)"
	}
	// Set the grid title.
	mp.Grid.SetTitle(fmt.Sprintf("%s, [lightgreen]%s!", welcome, username))

	// Fill in the table for the logged user.
	mp.SetUpLoggedTable(pages, 0)
}

// SetUpLoggedTable : Readies the MangaListTable to show information for a logged user.
// This includes the manga title and the publication status for the manga.
func (mp *MainPage) SetUpLoggedTable(pages *tview.Pages, offset int) {
	mangaTitleHeader := tview.NewTableCell("Title").
		SetAlign(tview.AlignCenter).
		SetTextColor(g.LoggedMainPageTitleColor).
		SetSelectable(false)
	pubStatusHeader := tview.NewTableCell("Pub. Status").
		SetAlign(tview.AlignLeft).
		SetTextColor(g.LoggedMainPagePubStatusColor).
		SetSelectable(false)
	mp.MangaListTable.SetCell(0, 0, mangaTitleHeader).
		SetCell(0, 1, pubStatusHeader).
		SetFixed(1, 0)

	go func() {
		// Get list of user's followed manga.
		mangaList, err := g.Dex.GetUserFollowedMangaList(g.OffsetRange, offset)
		if err != nil {
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				OKModal(pages, g.GenericAPIErrorModalID, "Error loading your followed manga.")
			})
			return
		} else if len(mangaList.Results) == 0 {
			noResCell := tview.NewTableCell("You have no followed manga!").SetSelectable(false)
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				mp.MangaListTable.SetCell(1, 0, noResCell)
			})
			return
		}

		// Get pagination numbers and fill table title.
		page, first, last := calculatePaginationData(offset, mangaList.Total)
		mp.MangaListTable.SetTitle(fmt.Sprintf("Your followed manga. Page %d (%d-%d).", page, first, last))

		// Set up input captures for the table.
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		SetMainPageTableHandlers(cancel, pages, mp, mangaList, offset, "")

		// Fill in each manga info row
		for i, mr := range mangaList.Results {
			select {
			case <-ctx.Done():
				return
			default:
				// Manga title cell.
				mtCell := tview.NewTableCell(fmt.Sprintf("%-50s", mr.Data.Attributes.Title["en"])).
					SetMaxWidth(50).SetTextColor(g.LoggedMainPageTitleColor)

				// Pub status cell.
				status := "-"
				if mr.Data.Attributes.Status != nil {
					status = strings.Title(*mr.Data.Attributes.Status)
				}
				sCell := tview.NewTableCell(fmt.Sprintf("%-15s", status)).
					SetMaxWidth(15).SetTextColor(g.LoggedMainPagePubStatusColor)

				g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
					mp.MangaListTable.SetCell(i+1, 0, mtCell).SetCell(i+1, 1, sCell)
				})
			}
		}
	}()
}

// SetUpGenericPage : Readies the MainPage to show a generic main page.
// You will be able to customise the grid and table title.
func (mp *MainPage) SetUpGenericPage(pages *tview.Pages, gridTitle, tableTitle string) {
	mp.Grid.SetTitle(gridTitle)

	// Fill in the generic manga list table.
	mp.SetUpGenericTable(pages, tableTitle, 0, "")
}

// SetUpGenericTable : Readies the MangaListTable to show generic manga information for a user.
// This includes the manga title, description, and tags.
func (mp *MainPage) SetUpGenericTable(pages *tview.Pages, tableTitle string, offset int, searchTitle string) {
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
	mp.MangaListTable.SetCell(0, 0, mangaTitleHeader).
		SetCell(0, 1, descHeader).
		SetCell(0, 2, tagHeader).
		SetFixed(1, 0)

	go func() {
		// Get manga list (search).
		// Set up search parameters
		params := url.Values{}
		params.Set("limit", strconv.Itoa(g.OffsetRange))
		params.Set("offset", strconv.Itoa(offset))
		if searchTitle != "" {
			params.Set("title", searchTitle)
		}
		mangaList, err := g.Dex.MangaList(params)
		if err != nil {
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				OKModal(pages, g.GenericAPIErrorModalID, "Error loading manga list.")
			})
			return
		} else if len(mangaList.Results) == 0 {
			noResCell := tview.NewTableCell("No results.").SetSelectable(false)
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				mp.MangaListTable.SetCell(1, 0, noResCell)
			})
			return
		}

		// Get pagination numbers and fill table title.
		page, first, last := calculatePaginationData(offset, mangaList.Total)
		mp.MangaListTable.SetTitle(fmt.Sprintf("%s Page %d (%d-%d).", tableTitle, page, first, last))

		// Set up input captures for the table.
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		SetMainPageTableHandlers(cancel, pages, mp, mangaList, offset, searchTitle)

		// Fill in each manga info row
		for i, mr := range mangaList.Results {
			select {
			case <-ctx.Done():
				return
			default:
				// Manga title cell.
				mtCell := tview.NewTableCell(fmt.Sprintf("%-40s", mr.Data.Attributes.Title["en"])).
					SetMaxWidth(40).SetTextColor(g.GuestMainPageTitleColor)

				// Description cell.
				desc := tview.Escape(fmt.Sprintf("%-70s", mr.Data.Attributes.Description["en"]))
				descCell := tview.NewTableCell(desc).SetMaxWidth(70).SetTextColor(g.GuestMainPageDescColor)

				// Tag cell.
				tags := make([]string, len(mr.Data.Attributes.Tags))
				for ti, tag := range mr.Data.Attributes.Tags {
					tags[ti] = tag.Attributes.Name["en"]
				}
				tagCell := tview.NewTableCell(strings.Join(tags, ", ")).SetTextColor(g.GuestMainPageTagColor)

				g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
					mp.MangaListTable.SetCell(i+1, 0, mtCell).
						SetCell(i+1, 1, descCell).
						SetCell(i+1, 2, tagCell)
				})
			}
		}
	}()
}

// calculatePaginationData : Calculates the current page and first/last entry number.
func calculatePaginationData(offset, total int) (int, int, int) {
	page := offset/g.OffsetRange + 1
	firstEntry := offset + 1
	lastEntry := page * g.OffsetRange

	if lastEntry > total {
		lastEntry = total
	}
	if firstEntry > lastEntry {
		firstEntry = lastEntry
	}

	return page, firstEntry, lastEntry
}
