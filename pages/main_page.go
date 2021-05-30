package pages

/*
Main Page shows the interface when the user first enters the application.

There are 2 main different interfaces - one for logged users and another for non-logged users.
Take note that the non-logged interface is also used by the Search Page.
*/

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

type MainPage struct {
	Grid           *tview.Grid  // The page grid.
	MangaListTable *tview.Table // The table that contains the list of manga.
	LoggedPage     bool         // To track whether the page is for logged user or not.
	CurrentOffset  int
	MaxOffset      int
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
		CurrentOffset:  0,
		MaxOffset:      0,
	}

	// Decide what to show for the main page.
	if g.Dex.IsLoggedIn() {
		mainPage.LoggedPage = true
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
	mp.SetUpLoggedTable(pages)
}

// SetUpLoggedTable : Readies the MangaListTable to show information for a logged user.
// This includes the manga title and the publication status for the manga.
func (mp *MainPage) SetUpLoggedTable(pages *tview.Pages) {
	// This function will always clear the table and selected function before drawing again.
	// Required as this page has pagination ability.
	mp.MangaListTable.Clear()
	mp.MangaListTable.SetSelectedFunc(func(row, column int) {})

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

	// Get pagination numbers and fill table title.
	// This is filler until the title is updated using the API response later in the goroutine.
	// Helps the user keep track of the current page when flipping pages very fast.
	page, first, last := mp.CalculatePaginationData()
	mp.MangaListTable.SetTitle(fmt.Sprintf("Your followed manga. Page %d (%d-%d). [::bu]Loading...", page, first, last))

	// Use context to stop goroutines that are no longer needed.
	// The page handler will induce cancel whenever the user flips pages.
	ctx, cancel := context.WithCancel(context.Background())
	SetMainPageTableHandlers(cancel, pages, mp, "")

	// Fill in each manga info row
	go func() {
		// NOTE: Due to pagination behaviour and API rate limits,
		// it is necessary to add a sleep here to allow the context to have enough time to react to cancels.
		// This is due to how users are able to hold down Ctrl+F/Ctrl+B respectively to switch pages very fast.
		time.Sleep(time.Millisecond * 40)
		defer cancel()

		// Get list of user's followed manga.
		var (
			mangaList *mangodex.MangaList
			err       error
		)
		select {
		case <-ctx.Done():
			return
		default:
			mangaList, err = g.Dex.GetUserFollowedMangaList(g.OffsetRange, mp.CurrentOffset)
			if err != nil {
				g.App.QueueUpdateDraw(func() {
					OKModal(pages, g.GenericAPIErrorModalID, "Error loading followed manga.")
				})
				return
			} else if len(mangaList.Results) == 0 {
				noResCell := tview.NewTableCell("You have no followed manga!").SetSelectable(false)
				g.App.QueueUpdateDraw(func() {
					mp.MangaListTable.SetCell(1, 0, noResCell)
					mp.MangaListTable.SetTitle(fmt.Sprintf("Your followed manga. Page %d (%d-%d).", page, first, last))
				})
				return
			}
			mp.MaxOffset = mangaList.Total // Note how the max offset is updated here.
		}

		// Get pagination numbers and fill table title.
		page, first, last = mp.CalculatePaginationData()
		g.App.QueueUpdateDraw(func() {
			mp.MangaListTable.SetTitle(fmt.Sprintf("Your followed manga. Page %d (%d-%d).", page, first, last))
		})

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

				g.App.QueueUpdateDraw(func() {
					mp.MangaListTable.SetCell(i+1, 0, mtCell).SetCell(i+1, 1, sCell)
				})

				// Set selected function for the table.
				// When user presses ENTER on a manga row, they are redirected to the manga page.
				// It is inside the for loop so user can press enter the moment they see an entry.
				mp.MangaListTable.SetSelectedFunc(func(row, column int) {
					// We do not need to worry about index out-of-range as we checked results earlier.
					ShowMangaPage(pages, &(mangaList.Results[row-1]))
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
	mp.SetUpGenericTable(pages, tableTitle, "")
}

// SetUpGenericTable : Readies the MangaListTable to show generic manga information for a user.
// This includes the manga title, description, and tags.
func (mp *MainPage) SetUpGenericTable(pages *tview.Pages, tableTitle string, searchTitle string) {
	// This function will always clear the table and selected function before drawing again.
	// Required as this page has pagination ability.
	mp.MangaListTable.Clear()
	mp.MangaListTable.SetSelectedFunc(func(row, column int) {})

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

	// Get pagination numbers and fill table title.
	// This is filler until the title is updated using the API response later in the goroutine.
	// Helps the user keep track of the current page when flipping pages very fast.
	page, first, last := mp.CalculatePaginationData()
	mp.MangaListTable.SetTitle(fmt.Sprintf("%s Page %d (%d-%d). [::bu]Loading...", tableTitle, page, first, last))

	// Use context to stop goroutines that are no longer needed.
	// The page handler will induce cancel whenever the user flips pages.
	ctx, cancel := context.WithCancel(context.Background())
	SetMainPageTableHandlers(cancel, pages, mp, searchTitle)

	// Fill in each manga info row
	go func() {
		// NOTE: Due to pagination behaviour and API rate limits,
		// it is necessary to add a sleep here to allow the context to have enough time to react to cancels.
		// This is due to how users are able to hold down Ctrl+F/Ctrl+B respectively to switch pages very fast.
		time.Sleep(time.Millisecond * 40)
		defer cancel()

		// Get manga list (search).
		// Set up search parameters
		params := url.Values{}
		params.Set("limit", strconv.Itoa(g.OffsetRange))
		params.Set("offset", strconv.Itoa(mp.CurrentOffset))
		if searchTitle != "" {
			params.Set("title", searchTitle)
		}
		var (
			mangaList *mangodex.MangaList
			err       error
		)
		select {
		case <-ctx.Done():
			return
		default:
			mangaList, err = g.Dex.MangaList(params)
			if err != nil {
				g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
					OKModal(pages, g.GenericAPIErrorModalID, "Error loading manga list.")
				})
				return
			} else if len(mangaList.Results) == 0 {
				noResCell := tview.NewTableCell("No results.").SetSelectable(false)
				g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
					mp.MangaListTable.SetCell(1, 0, noResCell)
					mp.MangaListTable.SetTitle(fmt.Sprintf("%s Page %d (%d-%d).", tableTitle, page, first, last))
				})
				return
			}
			mp.MaxOffset = mangaList.Total // Note how the max offset is updated here.
		}

		// Get pagination numbers and fill table title.
		page, first, last = mp.CalculatePaginationData()
		g.App.QueueUpdateDraw(func() {
			mp.MangaListTable.SetTitle(fmt.Sprintf("%s Page %d (%d-%d).", tableTitle, page, first, last))
		})

		for i, mr := range mangaList.Results {
			select {
			case <-ctx.Done():
				return
			default:
				// Manga title cell.
				mtCell := tview.NewTableCell(fmt.Sprintf("%-40s", mr.Data.Attributes.Title["en"])).
					SetMaxWidth(40).SetTextColor(g.GuestMainPageTitleColor)

				// Description cell. Truncate description to improve loading times.
				desc := tview.Escape(fmt.Sprintf("%-60s",
					strings.SplitN(tview.Escape(mr.Data.Attributes.Description["en"]), "\n", 2)[0]))
				descCell := tview.NewTableCell(desc).SetMaxWidth(60).SetTextColor(g.GuestMainPageDescColor)

				// Tag cell.
				tags := make([]string, len(mr.Data.Attributes.Tags))
				for ti, tag := range mr.Data.Attributes.Tags {
					tags[ti] = tag.Attributes.Name["en"]
				}
				tagCell := tview.NewTableCell(strings.Join(tags, ", ")).SetTextColor(g.GuestMainPageTagColor)

				g.App.QueueUpdateDraw(func() {
					mp.MangaListTable.SetCell(i+1, 0, mtCell).
						SetCell(i+1, 1, descCell).
						SetCell(i+1, 2, tagCell)
				})

				// Set selected function for the table.
				// When user presses ENTER on a manga row, they are redirected to the manga page.
				// It is inside the for loop so user can press enter the moment they see an entry.
				mp.MangaListTable.SetSelectedFunc(func(row, column int) {
					// We do not need to worry about index out-of-range as we checked results earlier.
					ShowMangaPage(pages, &(mangaList.Results[row-1]))
				})
			}
		}
	}()
}

// CalculatePaginationData : Calculates the current page and first/last entry number.
func (mp *MainPage) CalculatePaginationData() (int, int, int) {
	page := mp.CurrentOffset/g.OffsetRange + 1
	firstEntry := mp.CurrentOffset + 1
	lastEntry := page * g.OffsetRange

	if lastEntry > mp.MaxOffset {
		lastEntry = mp.MaxOffset
	}
	if firstEntry > lastEntry {
		firstEntry = lastEntry
	}

	return page, firstEntry, lastEntry
}
