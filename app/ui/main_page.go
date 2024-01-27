package ui

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/darylhjd/mangadesk/app/core"
	"github.com/darylhjd/mangadesk/app/ui/utils"
	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
)

const (
	offsetRange = 100
	loadDelay   = time.Millisecond * 50
	maxOffset   = 10000
)

// MainPage : This struct contains the grid and the entry table.
// In addition, it also keeps track of whether to show followed/popular manga based on login status
// as well as the entry offset.
type MainPage struct {
	Grid          *tview.Grid  // The page grid.
	Table         *tview.Table // The table that contains the list of manga.
	CurrentOffset int
	MaxOffset     int

	cWrap *utils.ContextWrapper // For context cancellation.
}

// ShowMainPage : Make the app show the main page.
func ShowMainPage() {
	// Create the new main page
	log.Println("Creating new main page...")
	mainPage := newMainPage()

	core.App.TView.SetFocus(mainPage.Grid)
	core.App.PageHolder.AddAndSwitchToPage(utils.MainPageID, mainPage.Grid, true)
}

// newMainPage : Creates a new main page.
func newMainPage() *MainPage {
	var dimensions []int
	for i := 0; i < 15; i++ {
		dimensions = append(dimensions, -1)
	}
	grid := utils.NewGrid(dimensions, dimensions)
	// Set grid attributes.
	grid.SetTitleColor(utils.MainPageGridTitleColor).
		SetBorderColor(utils.MainPageGridBorderColor).
		SetBorder(true)

	// Create the base main table.
	table := tview.NewTable()
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(utils.MainPageTableBorderColor).
		SetTitleColor(utils.MainPageTableTitleColor).
		SetBorder(true)

	// Add the table to the grid. Table spans the whole page.
	grid.AddItem(table, 0, 0, 15, 15, 0, 0, true)

	ctx, cancel := context.WithCancel(context.Background())
	mainPage := &MainPage{
		Grid:  grid,
		Table: table,
		cWrap: &utils.ContextWrapper{
			Ctx:    ctx,
			Cancel: cancel,
		},
	}

	// Check what kind of main page to show to the user.
	if core.App.Client.Auth.IsLoggedIn() {
		mainPage.setLogged()
	} else {
		mainPage.setGuest()
	}
	return mainPage
}

// setLogged : Set up the MainPage for a logged user.
func (p *MainPage) setLogged() {
	log.Println("Using logged main page.")
	go p.setLoggedGrid()
	go p.setLoggedTable()
}

// setLoggedGrid : Show logged grid title.
func (p *MainPage) setLoggedGrid() {
	log.Println("Setting logged grid...")
	var username string
	if u, err := core.App.Client.User.GetLoggedUser(); err != nil {
		log.Println(fmt.Sprintf("Error getting user info: %s", err.Error()))
	} else {
		username = u.Data.Attributes.Username
	}

	core.App.TView.QueueUpdateDraw(func() {
		p.Grid.SetTitle(fmt.Sprintf("Welcome to MangaDex, [lightgreen]%s!", username))
	})
	log.Println("Finished setting logged grid.")
}

// setLoggedTable : Show logged table items and title.
func (p *MainPage) setLoggedTable() {
	log.Println("Setting logged table...")
	ctx, cancel := p.cWrap.ResetContext()
	// Set handlers
	p.setHandlers(cancel, nil)

	time.Sleep(loadDelay)
	defer cancel()

	core.App.TView.QueueUpdateDraw(func() {
		// Clear current entries.
		p.Table.Clear()

		// Set headers.
		titleHeader := tview.NewTableCell("Title").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.LoggedMainPageTitleColor).
			SetSelectable(false)
		pubStatusHeader := tview.NewTableCell("Pub. Status").
			SetAlign(tview.AlignLeft).
			SetTextColor(utils.LoggedMainPagePubStatusColor).
			SetSelectable(false)
		p.Table.SetCell(0, 0, titleHeader).
			SetCell(0, 1, pubStatusHeader).
			SetFixed(1, 0)

		// Set table title.
		page, first, last := p.calculatePaginationData()
		p.Table.SetTitle(fmt.Sprintf("Followed manga. Page %d (%d-%d). [::bu]Loading...", page, first, last))
	})

	// Get the list of the user's followed manga.
	if p.cWrap.ToCancel(ctx) {
		return
	}
	followed, err := core.App.Client.User.GetUserFollowedMangaList(
		offsetRange, p.CurrentOffset, []string{mangodex.AuthorRel})
	if err != nil {
		log.Printf("Error getting followed manga: %s\n", err.Error())
		core.App.TView.QueueUpdateDraw(func() {
			modal := okModal(utils.GenericAPIErrorModalID, "Error getting followed manga.\nCheck logs for details.")
			ShowModal(utils.GenericAPIErrorModalID, modal)
		})
		return
	}

	// Sort the list based on manga title.
	sort.Slice(followed.Data, func(i, j int) bool {
		return followed.Data[i].GetTitle("en") < followed.Data[j].GetTitle("en")
	})

	// Update offset details.
	p.MaxOffset = int(math.Min(float64(followed.Total), maxOffset))

	// Show followed manga.
	if p.MaxOffset == 0 {
		core.App.TView.QueueUpdateDraw(func() {
			noResCell := tview.NewTableCell("You have no followed manga!").SetSelectable(false)
			p.Table.SetCell(1, 0, noResCell)
		})
		return
	}

	// Update table title.
	page, first, last := p.calculatePaginationData()
	core.App.TView.QueueUpdateDraw(func() {
		p.Table.SetTitle(fmt.Sprintf("Followed manga. Page %d (%d-%d).", page, first, last))
	})

	// Fill in the details
	for index := 0; index < len(followed.Data); index++ {
		if p.cWrap.ToCancel(ctx) {
			return
		}
		manga := followed.Data[index]
		// Set title and publishing status cells.
		// Title
		mtCell := tview.NewTableCell(fmt.Sprintf("%-50s", manga.GetTitle("en"))).
			SetMaxWidth(50).SetTextColor(utils.LoggedMainPageTitleColor).SetReference(&manga)

		// Publishing Status.
		sCell := tview.NewTableCell(strings.Title(fmt.Sprintf("%-15s", *manga.Attributes.Status))).
			SetMaxWidth(15).SetTextColor(utils.LoggedMainPagePubStatusColor)

		p.Table.SetCell(index+1, 0, mtCell).SetCell(index+1, 1, sCell)
	}
	core.App.TView.QueueUpdateDraw(func() {
		p.Table.Select(1, 0)
		p.Table.ScrollToBeginning()
	})
	log.Println("Finished setting logged table.")
}

// setGuest : Set up the main page for a guest user.
func (p *MainPage) setGuest() {
	log.Println("Using guest main page.")
	go p.setGuestGrid()
	go p.setGuestTable(nil)
}

// setGuestGrid : Show guest grid title.
func (p *MainPage) setGuestGrid() {
	log.Println("Setting guest grid...")
	core.App.TView.QueueUpdateDraw(func() {
		p.Grid.SetTitle("Welcome to MangaDex, [yellow]Guest!")
	})
	log.Println("Finished setting guest grid.")
}

// setGuestTable : Show guest table items and title. This function is also used to create a search table.
// Whether we are setting up the table for the guest main page or a search page depends on whether
// searchParams is nil. If it is nil, then it is not a search, otherwise we are searching.
func (p *MainPage) setGuestTable(searchParams *SearchParams) {
	log.Println("Setting guest table...")
	ctx, cancel := p.cWrap.ResetContext()
	// Set the handlers
	p.setHandlers(cancel, searchParams)

	time.Sleep(loadDelay)
	defer cancel()

	tableTitle := "Popular manga"
	if searchParams != nil {
		tableTitle = "Search Results"
	}

	core.App.TView.QueueUpdateDraw(func() {
		// Clear current entries
		p.Table.Clear()

		// Set headers.
		titleHeader := tview.NewTableCell("Manga").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageTitleColor).
			SetSelectable(false)
		descHeader := tview.NewTableCell("Description").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageDescColor).
			SetSelectable(false)
		tagHeader := tview.NewTableCell("Tags").
			SetAlign(tview.AlignCenter).
			SetTextColor(utils.GuestMainPageTagColor).
			SetSelectable(false)
		p.Table.SetCell(0, 0, titleHeader).
			SetCell(0, 1, descHeader).
			SetCell(0, 2, tagHeader).
			SetFixed(1, 0)

		// Set table title.
		page, first, last := p.calculatePaginationData()
		p.Table.SetTitle(fmt.Sprintf("%s. Page %d (%d-%d). [::bu]Loading...", tableTitle, page, first, last))
	})

	// Get list of manga.
	params := p.setGuestSearchParams(searchParams)
	if p.cWrap.ToCancel(ctx) {
		return
	}
	list, err := core.App.Client.Manga.GetMangaList(*params)
	if err != nil {
		log.Println(err.Error())
		core.App.TView.QueueUpdateDraw(func() {
			modal := okModal(utils.GenericAPIErrorModalID, "Error getting manga list.\nCheck logs for details.")
			ShowModal(utils.GenericAPIErrorModalID, modal)
		})
		return
	}

	// Update offset details.
	p.MaxOffset = int(math.Min(float64(list.Total), maxOffset))

	// Show followed manga.
	if p.MaxOffset == 0 {
		core.App.TView.QueueUpdateDraw(func() {
			noResCell := tview.NewTableCell("No results!").SetSelectable(false)
			p.Table.SetCell(1, 0, noResCell)
		})
		return
	}

	// Update table title.
	page, first, last := p.calculatePaginationData()
	core.App.TView.QueueUpdateDraw(func() {
		p.Table.SetTitle(fmt.Sprintf("%s. Page %d (%d-%d).", tableTitle, page, first, last))
	})

	// Fill in the details
	for index := 0; index < len(list.Data); index++ {
		if p.cWrap.ToCancel(ctx) {
			return
		}
		manga := list.Data[index]
		// Manga title cell.
		mtCell := tview.NewTableCell(fmt.Sprintf("%-40s", manga.GetTitle("en"))).
			SetMaxWidth(40).SetTextColor(utils.GuestMainPageTitleColor).SetReference(&manga)

		// Description cell. Truncate description to improve loading times.
		desc := tview.Escape(fmt.Sprintf("%-60s",
			strings.SplitN(tview.Escape(manga.GetDescription("en")), "\n", 2)[0]))
		descCell := tview.NewTableCell(desc).SetMaxWidth(60).SetTextColor(utils.GuestMainPageDescColor)

		// Tag cell.
		tags := make([]string, len(manga.Attributes.Tags))
		for i, tag := range manga.Attributes.Tags {
			tags[i] = tag.GetName("en")
		}
		tagCell := tview.NewTableCell(strings.Join(tags, ", ")).SetTextColor(utils.GuestMainPageTagColor)

		p.Table.SetCell(index+1, 0, mtCell).
			SetCell(index+1, 1, descCell).
			SetCell(index+1, 2, tagCell)
	}
	core.App.TView.QueueUpdateDraw(func() {
		p.Table.Select(1, 0)
		p.Table.ScrollToBeginning()
	})
	log.Println("Finished setting guest table.")
}

// setGuestSearchParams : Helper function to set up query parameters for the guest table.
func (p *MainPage) setGuestSearchParams(searchParams *SearchParams) *url.Values {
	// Set up offset parameters
	params := url.Values{}

	// Set limits and offset
	params.Set("limit", strconv.Itoa(offsetRange))
	params.Set("offset", strconv.Itoa(p.CurrentOffset))

	// Set content ratings
	ratings := []string{mangodex.Safe, mangodex.Suggestive, mangodex.Erotica}
	switch searchParams != nil {
	case true: // If it is a search, then we follow settings for this search.
		if !searchParams.explicit {
			ratings = append(ratings, mangodex.Porn)
		}
	case false: // Else, we follow the configuration settings set by the user.
		if core.App.Config.ExplicitContent {
			ratings = append(ratings, mangodex.Porn)
		}
	}
	for _, rating := range ratings {
		params.Add("contentRating[]", rating)
	}

	// Set order of results.
	if searchParams != nil { // If fpr searching, we sort by relevancy.
		params.Set("order[relevance]", "desc")
	} else { // Else, we sort by popular manga (based on follow count)
		params.Set("order[followedCount]", "desc")
	}

	// Set search term (if for search)
	if searchParams != nil {
		log.Printf("Settings guest table for search: \"%s\"\n", searchParams.term)
		params.Set("title", searchParams.term)
	}

	// Include Author relationship
	params.Set("includes[]", mangodex.AuthorRel)
	return &params
}

// calculatePaginationData : Calculates the current page and first/last entry number.
// Returns (pageNo, firstEntry, lastEntry).
func (p *MainPage) calculatePaginationData() (int, int, int) {
	page := p.CurrentOffset/offsetRange + 1
	firstEntry := p.CurrentOffset + 1
	lastEntry := page * offsetRange

	if lastEntry > p.MaxOffset {
		lastEntry = p.MaxOffset
	}
	if firstEntry > lastEntry {
		firstEntry = lastEntry
	}

	return page, firstEntry, lastEntry
}
