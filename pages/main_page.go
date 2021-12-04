package pages

import (
	"fmt"
	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/darylhjd/mangadesk/core"
)

// MainPage : This struct contains the grid and the entry table.
// In addition, it also keeps track of whether to show followed/popular manga based on login status
// as well as the entry offset.
type MainPage struct {
	Grid          *tview.Grid  // The page grid.
	Table         *tview.Table // The table that contains the list of manga.
	LoggedPage    bool         // To track whether the page is for logged user or not.
	CurrentOffset int
	MaxOffset     int
}

// NewMainPage : Creates a new main page.
func NewMainPage() *MainPage {
	var dimensions []int
	for i := 0; i < 15; i++ {
		dimensions = append(dimensions, -1)
	}
	grid := NewGrid(dimensions, dimensions)
	// Set grid attributes.
	grid.SetTitleColor(MainPageGridTitleColor).
		SetBorderColor(MainPageGridBorderColor).
		SetBorder(true)

	// Create the base main table.
	table := tview.NewTable()
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(MainPageTableBorderColor).
		SetTitleColor(MainPageTableTitleColor).
		SetBorder(true)

	// Add the table to the grid. Table spans the whole page.
	grid.AddItem(table, 0, 0, 15, 15, 0, 0, true)

	mainPage := &MainPage{
		Grid:  grid,
		Table: table,
	}

	// Check what kind of main page to show to the user.
	if core.App.Client.Auth.IsLoggedIn() {
		mainPage.LoggedPage = true
		mainPage.SetLogged()
	} else {
		mainPage.SetGuest()
	}
	return mainPage
}

// SetLogged : Set up the MainPage for a logged user.
func (p *MainPage) SetLogged() {
	go p.setLoggedGrid()
	go p.setLoggedTable()
}

// setLoggedGrid : Show logged grid title.
func (p *MainPage) setLoggedGrid() {
	var username string
	if u, err := core.App.Client.User.GetLoggedUser(); err != nil {
		log.Println(fmt.Sprintf("Error getting user info: %s", err.Error()))
	} else {
		username = u.Data.Attributes.Username
	}

	core.App.ViewApp.QueueUpdateDraw(func() {
		p.Grid.SetTitle(fmt.Sprintf("Welcome to MangaDex, [lightgreen]%s!", username))
	})
}

// setLoggedTable : Show logged table items and title.
func (p *MainPage) setLoggedTable() {
	core.App.ViewApp.QueueUpdateDraw(func() {
		// Clear current entries.
		p.Table.Clear()

		// Set headers.
		titleHeader := tview.NewTableCell("Title").
			SetAlign(tview.AlignCenter).
			SetTextColor(LoggedMainPageTitleColor).
			SetSelectable(false)
		pubStatusHeader := tview.NewTableCell("Pub. Status").
			SetAlign(tview.AlignLeft).
			SetTextColor(LoggedMainPagePubStatusColor).
			SetSelectable(false)
		p.Table.SetCell(0, 0, titleHeader).
			SetCell(0, 1, pubStatusHeader).
			SetFixed(1, 0)

		// Set table title.
		page, first, last := p.CalculatePaginationData()
		p.Table.SetTitle(fmt.Sprintf("Followed manga. Page %d (%d-%d). [::bu]Loading...", page, first, last))
	})

	// Get the list of the user's followed manga.
	followed, err := core.App.Client.User.GetUserFollowedMangaList(
		OffsetRange, p.CurrentOffset, []string{mangodex.AuthorRel})
	if err != nil {
		log.Println(err.Error())
		modal := OKModal(GenericAPIErrorModalID, "Error getting followed manga.\nCheck logs for details.")
		core.App.ShowModal(GenericAPIErrorModalID, modal)
		return
	}

	// Update offset details.
	p.MaxOffset = followed.Total

	// Show followed manga.
	if p.MaxOffset == 0 {
		core.App.ViewApp.QueueUpdateDraw(func() {
			noResCell := tview.NewTableCell("You have no followed manga!").SetSelectable(false)
			p.Table.SetCell(1, 0, noResCell)
		})
		return
	}

	// Update table title.
	page, first, last := p.CalculatePaginationData()
	core.App.ViewApp.QueueUpdateDraw(func() {
		p.Table.SetTitle(fmt.Sprintf("Followed manga. Page %d (%d-%d).", page, first, last))
	})

	p.Table.SetSelectedFunc(func(row, _ int) {
		core.App.ShowMangaPage((p.Table.GetCell(row, 0).GetReference()).(*mangodex.Manga))
	})

	// Fill in the details
	for index, manga := range followed.Data {
		// Set title and publishing status cells.
		// Title
		mtCell := tview.NewTableCell(fmt.Sprintf("%-50s", manga.GetTitle("en"))).
			SetMaxWidth(50).SetTextColor(LoggedMainPageTitleColor).SetReference(&manga)

		// Publishing Status.
		sCell := tview.NewTableCell(fmt.Sprintf("%-15s", *manga.Attributes.Status)).
			SetMaxWidth(15).SetTextColor(LoggedMainPagePubStatusColor)

		core.App.ViewApp.QueueUpdateDraw(func() {
			p.Table.SetCell(index+1, 0, mtCell).SetCell(index+1, 1, sCell)
		})
	}
}

// SetGuest : Set up the main page for a guest user.
func (p *MainPage) SetGuest() {
	go p.setGuestGrid()
	go p.setGuestTable()
}

// setGuestGrid : Show guest grid title.
func (p *MainPage) setGuestGrid() {
	core.App.ViewApp.QueueUpdateDraw(func() {
		p.Grid.SetTitle("Welcome to MangaDex, [yellow]Guest!")
	})
}

// setGuestTable : Show guest table items and title.
func (p *MainPage) setGuestTable() {
	core.App.ViewApp.QueueUpdateDraw(func() {
		// Clear current entries
		p.Table.Clear()

		// Set headers.
		titleHeader := tview.NewTableCell("Manga").
			SetAlign(tview.AlignCenter).
			SetTextColor(GuestMainPageTitleColor).
			SetSelectable(false)
		descHeader := tview.NewTableCell("Description").
			SetAlign(tview.AlignCenter).
			SetTextColor(GuestMainPageDescColor).
			SetSelectable(false)
		tagHeader := tview.NewTableCell("Tags").
			SetAlign(tview.AlignCenter).
			SetTextColor(GuestMainPageTagColor).
			SetSelectable(false)
		p.Table.SetCell(0, 0, titleHeader).
			SetCell(0, 1, descHeader).
			SetCell(0, 2, tagHeader).
			SetFixed(1, 0)

		// Set table title.
		// Set table title.
		page, first, last := p.CalculatePaginationData()
		p.Table.SetTitle(fmt.Sprintf("Popular manga. Page %d (%d-%d). [::bu]Loading...", page, first, last))
	})

	// Get list of manga.
	// Set up search parameters
	params := url.Values{}
	params.Set("limit", strconv.Itoa(OffsetRange))
	params.Set("offset", strconv.Itoa(p.CurrentOffset))
	// If user wants explicit content.
	ratings := []string{"safe", "suggestive", "erotica"}
	if core.App.Config.ExplicitContent {
		ratings = append(ratings, "pornographic")
	}
	for _, rating := range ratings {
		params.Add("contentRating[]", rating)
	}
	// Include Author relationship
	params.Set("includes[]", mangodex.AuthorRel)

	list, err := core.App.Client.Manga.GetMangaList(params)
	if err != nil {
		log.Println(err.Error())
		modal := OKModal(GenericAPIErrorModalID, "Error getting manga list.\nCheck logs for details.")
		core.App.ShowModal(GenericAPIErrorModalID, modal)
		return
	}

	// Update offset details.
	p.MaxOffset = list.Total

	// Show followed manga.
	if p.MaxOffset == 0 {
		core.App.ViewApp.QueueUpdateDraw(func() {
			noResCell := tview.NewTableCell("No manga entry!").SetSelectable(false)
			p.Table.SetCell(1, 0, noResCell)
		})
		return
	}

	// Update table title.
	page, first, last := p.CalculatePaginationData()
	core.App.ViewApp.QueueUpdateDraw(func() {
		p.Table.SetTitle(fmt.Sprintf("Popular manga. Page %d (%d-%d).", page, first, last))
	})

	p.Table.SetSelectedFunc(func(row, _ int) {
		core.App.ShowMangaPage((p.Table.GetCell(row, 0).GetReference()).(*mangodex.Manga))
	})

	// Fill in the details
	for index, manga := range list.Data {
		// Manga title cell.
		mtCell := tview.NewTableCell(fmt.Sprintf("%-40s", manga.GetTitle("en"))).
			SetMaxWidth(40).SetTextColor(GuestMainPageTitleColor).SetReference(&manga)

		// Description cell. Truncate description to improve loading times.
		desc := tview.Escape(fmt.Sprintf("%-60s",
			strings.SplitN(tview.Escape(manga.GetDescription("en")), "\n", 2)[0]))
		descCell := tview.NewTableCell(desc).SetMaxWidth(60).SetTextColor(GuestMainPageDescColor)

		// Tag cell.
		tags := make([]string, len(manga.Attributes.Tags))
		for i, tag := range manga.Attributes.Tags {
			tags[i] = tag.GetName("en")
		}
		tagCell := tview.NewTableCell(strings.Join(tags, ", ")).SetTextColor(GuestMainPageTagColor)

		core.App.ViewApp.QueueUpdateDraw(func() {
			p.Table.SetCell(index+1, 0, mtCell).
				SetCell(index+1, 1, descCell).
				SetCell(index+1, 2, tagCell)
		})
	}
}

// CalculatePaginationData : Calculates the current page and first/last entry number.
func (p *MainPage) CalculatePaginationData() (int, int, int) {
	page := p.CurrentOffset/OffsetRange + 1
	firstEntry := p.CurrentOffset + 1
	lastEntry := page * OffsetRange

	if lastEntry > p.MaxOffset {
		lastEntry = p.MaxOffset
	}
	if firstEntry > lastEntry {
		firstEntry = lastEntry
	}

	return page, firstEntry, lastEntry
}
