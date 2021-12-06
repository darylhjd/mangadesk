package ui

import (
	"fmt"
	"github.com/darylhjd/mangadesk/app/core"
	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"net/url"
	"os"
	"strconv"
	"sync"
)

// MangaPage : This struct contains the required primitives for the manga page.
type MangaPage struct {
	Manga *mangodex.Manga
	Grid  *tview.Grid
	Info  *tview.TextView
	Table *tview.Table

	Selected      []int // Keep track of which chapters have been selected by user.
	SelectedMutex *sync.Mutex
}

// ShowMangaPage : Make the app show the manga page.
func ShowMangaPage(manga *mangodex.Manga) {
	mangaPage := newMangaPage(manga)

	core.App.TView.SetFocus(mangaPage.Grid)
	core.App.PageHolder.AddAndSwitchToPage(MangaPageID, mangaPage.Grid, true)
}

// newMangaPage : Creates a new manga page.
func newMangaPage(manga *mangodex.Manga) *MangaPage {
	var dimensions []int
	for i := 0; i < 15; i++ {
		dimensions = append(dimensions, -1)
	}
	grid := newGrid(dimensions, dimensions)
	// Set grid attributes
	grid.SetTitleColor(MangaPageGridTitleColor).
		SetBorderColor(MangaPageGridBorderColor).
		SetTitle("Manga Information").
		SetBorder(true)

	// Use a TextView for basic information of the manga.
	info := tview.NewTextView()
	// Set textview attributes
	info.SetWrap(true).SetWordWrap(true).
		SetBorderColor(MangaPageInfoViewBorderColor).
		SetTitleColor(MangaPageInfoViewTitleColor).
		SetTitle("About").
		SetBorder(true)

	// Use a table to show the chapters for the manga.
	table := tview.NewTable()
	// Set chapter headers
	numHeader := tview.NewTableCell("Chap").
		SetTextColor(MangaPageChapNumColor).
		SetSelectable(false)
	titleHeader := tview.NewTableCell("Name").
		SetTextColor(MangaPageTitleColor).
		SetSelectable(false)
	downloadHeader := tview.NewTableCell("Download Status").
		SetTextColor(MangaPageDownloadStatColor).
		SetSelectable(false)
	scanGroupHeader := tview.NewTableCell("ScanGroup").
		SetTextColor(MangaPageScanGroupColor).
		SetSelectable(false)
	readMarkerHeader := tview.NewTableCell("Read Status").
		SetTextColor(MangaPageReadStatColor).
		SetSelectable(false)
	table.SetCell(0, 0, numHeader).
		SetCell(0, 1, titleHeader).
		SetCell(0, 2, downloadHeader).
		SetCell(0, 3, scanGroupHeader).
		SetCell(0, 4, readMarkerHeader).
		SetFixed(1, 0)
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(MangaPageTableBorderColor).
		SetTitle("Chapters").
		SetTitleColor(MangaPageTableTitleColor).
		SetBorder(true)

	// Add info and table to the grid. Set the focus to the chapter table.
	grid.AddItem(info, 0, 0, 5, 15, 0, 0, false).
		AddItem(table, 5, 0, 10, 15, 0, 0, true).
		AddItem(info, 0, 0, 15, 5, 0, 80, false).
		AddItem(table, 0, 5, 15, 10, 0, 80, true)

	mangaPage := &MangaPage{
		Manga:         manga,
		Grid:          grid,
		Info:          info,
		Table:         table,
		Selected:      []int{},
		SelectedMutex: &sync.Mutex{},
	}

	// Set up values
	go mangaPage.setMangaInfo()
	go mangaPage.setChapterTable()

	return mangaPage
}

// setMangaInfo : Set up manga information.
func (p *MangaPage) setMangaInfo() {
	// Title
	title := p.Manga.GetTitle("en")

	// Author
	var author string
	for _, relation := range p.Manga.Relationships {
		if relation.Type == mangodex.AuthorRel {
			author = relation.Attributes.(*mangodex.AuthorAttributes).Name
			break
		}
	}

	// Status
	status := *p.Manga.Attributes.Status

	// Description
	desc := tview.Escape(p.Manga.GetDescription("en"))

	// Set up information text.
	infoText := fmt.Sprintf("Title: %s\n\nAuthor: %s\nStatus: %s\n\nDescription:\n%s",
		title, author, status, desc)

	core.App.TView.QueueUpdateDraw(func() {
		p.Info.SetText(infoText)
	})
}

// setChapterTable : Fill up the chapter table.
func (p *MangaPage) setChapterTable() {
	// Show loading status so user knows it's loading.
	core.App.TView.QueueUpdateDraw(func() {
		loadingCell := tview.NewTableCell("Loading...").SetSelectable(false)
		p.Table.SetCell(1, 1, loadingCell)
	})

	// Get all chapters
	chapters, err := p.getAllChapters()
	if err != nil {
		log.Println(fmt.Sprintf("Error getting manga chapters: %s", err.Error()))
		modal := okModal(GenericAPIErrorModalID, "Error getting manga chapters.\nCheck log for details.")
		ShowModal(GenericAPIErrorModalID, modal)
		return
	}
	// Get the chapter read markers.
	var markers map[string]struct{}
	if core.App.Client.Auth.IsLoggedIn() {
		markerResponse, err := core.App.Client.Chapter.GetReadMangaChapters(p.Manga.ID)
		if err != nil {
			log.Println(fmt.Sprintf("Error getting chapter read markers: %s", err.Error()))
			modal := okModal(GenericAPIErrorModalID, "Error getting chapter read markers.\nCheck log for details.")
			ShowModal(GenericAPIErrorModalID, modal)
			return
		}
		for _, marker := range markerResponse.Data {
			markers[marker] = struct{}{}
		}
	}

	// TODO: Add manga page input handlers

	// Fill in the chapters
	for index, chapter := range chapters {
		// Chapter Number
		chapterNumCell := tview.NewTableCell(
			fmt.Sprintf("%-6s %s", chapter.GetChapterNum(), chapter.Attributes.TranslatedLanguage)).
			SetMaxWidth(10).SetTextColor(MangaPageChapNumColor).SetReference(&chapter)

		// Chapter title
		titleCell := tview.NewTableCell(fmt.Sprintf("%-30s", chapter.GetTitle())).SetMaxWidth(30).
			SetTextColor(MangaPageTitleColor)

		// Chapter download status
		var downloadStatus string
		// Check for the presence of the download folder.
		if _, err = os.Stat(p.getDownloadFolder(&chapter)); err == nil {
			downloadStatus = "Y"
		}
		downloadCell := tview.NewTableCell(downloadStatus).SetTextColor(MangaPageDownloadStatColor)

		// Scanlation group
		var scanGroup string
		for _, relation := range p.Manga.Relationships {
			if relation.Type == mangodex.ScanlationGroupRel {
				scanGroup = relation.Attributes.(*mangodex.ScanlationGroupAttributes).Name
				break
			}
		}
		scanGroupCell := tview.NewTableCell(fmt.Sprintf("%-15s", scanGroup)).SetMaxWidth(15).
			SetTextColor(MangaPageScanGroupColor)

		// Read marker
		var read string
		if !core.App.Client.Auth.IsLoggedIn() {
			read = "Not logged in!"
		} else if _, ok := markers[chapter.ID]; ok {
			read = "Y"
		}
		readCell := tview.NewTableCell(read).SetTextColor(MangaPageReadStatColor)

		core.App.TView.QueueUpdateDraw(func() {
			p.Table.SetCell(index+1, 0, chapterNumCell).
				SetCell(index+1, 1, titleCell).
				SetCell(index+1, 2, downloadCell).
				SetCell(index+1, 3, scanGroupCell)

			if !core.App.Client.Auth.IsLoggedIn() {
				if index == 0 {
					p.Table.SetCell(index+1, 4, readCell)
				}
			} else {
				p.Table.SetCell(index+1, 4, readCell)
			}
		})
	}
}

// getAllChapters : Get all chapters for the manga.
func (p *MangaPage) getAllChapters() ([]mangodex.Chapter, error) {
	// Set up query parameters.
	params := url.Values{}
	params.Set("limit", "500")
	// Get all chapters with user's specified languages
	for _, lang := range core.App.Config.Languages {
		params.Add("translatedLanguage[]", lang)
	}
	// Show the latest chapters first.
	params.Set("order[chapter]", "desc")
	// Show required explicit chapters
	ratings := []string{mangodex.Safe, mangodex.Suggestive, mangodex.Erotica}
	if p.Manga.Attributes.ContentRating != nil && *p.Manga.Attributes.ContentRating == mangodex.Porn {
		ratings = append(ratings, mangodex.Porn)
	}
	for _, rating := range ratings {
		params.Add("contentRating[]", rating)
	}
	// Also get the scanlation group for the chapter
	params.Add("includes[]", mangodex.ScanlationGroupRel)

	var (
		chapters   []mangodex.Chapter
		currOffset = 0
	)
	for {
		params.Set("offset", strconv.Itoa(currOffset))
		list, err := core.App.Client.Chapter.GetMangaChapters(p.Manga.ID, params)
		if err != nil {
			return []mangodex.Chapter{}, err
		}
		chapters = append(chapters, list.Data...)
		currOffset += 500
		if currOffset >= list.Total {
			break
		}
	}
	return chapters, nil
}

// MarkChapterSelected : Mark a chapter as being selected by the user on the main page table.
func (p *MangaPage) markChapterSelected(row int) {
	chapterCell := p.Table.GetCell(row, 0)
	chapterCell.SetTextColor(tcell.ColorBlack).SetBackgroundColor(MangaPageHighlightColor)
}

// MarkChapterUnselected : Mark a chapter as being unselected by the user on the main page table.
func (p *MangaPage) markChapterUnselected(row int) {
	chapterCell := p.Table.GetCell(row, 0)
	chapterCell.SetTextColor(MangaPageChapNumColor).SetBackgroundColor(tcell.ColorBlack)
}
