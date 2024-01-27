package ui

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/darylhjd/mangadesk/app/core"
	"github.com/darylhjd/mangadesk/app/ui/utils"
	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	chapterOffsetRange    = 500
	contextCancelledError = "CANCELLED"
	readStatus            = "Y"
)

// MangaPage : This struct contains the required primitives for the manga page.
type MangaPage struct {
	Manga *mangodex.Manga
	Grid  *tview.Grid
	Info  *tview.TextView
	Table *tview.Table

	sWrap *utils.SelectorWrapper
	cWrap *utils.ContextWrapper // For context cancellation.
}

// ShowMangaPage : Make the app show the manga page.
func ShowMangaPage(manga *mangodex.Manga) {
	mangaPage := newMangaPage(manga)

	core.App.TView.SetFocus(mangaPage.Grid)
	core.App.PageHolder.AddAndSwitchToPage(utils.MangaPageID, mangaPage.Grid, true)
}

// newMangaPage : Creates a new manga page.
func newMangaPage(manga *mangodex.Manga) *MangaPage {
	var dimensions []int
	for i := 0; i < 15; i++ {
		dimensions = append(dimensions, -1)
	}
	grid := utils.NewGrid(dimensions, dimensions)
	// Set grid attributes
	grid.SetTitleColor(utils.MangaPageGridTitleColor).
		SetBorderColor(utils.MangaPageGridBorderColor).
		SetTitle("Manga Information").
		SetBorder(true)

	// Use a TextView for basic information of the manga.
	info := tview.NewTextView()
	// Set textview attributes
	info.SetWrap(true).SetWordWrap(true).
		SetBorderColor(utils.MangaPageInfoViewBorderColor).
		SetTitleColor(utils.MangaPageInfoViewTitleColor).
		SetTitle("About").
		SetBorder(true)

	// Use a table to show the chapters for the manga.
	table := tview.NewTable()
	// Set chapter headers
	numHeader := tview.NewTableCell("Chap").
		SetTextColor(utils.MangaPageChapNumColor).
		SetSelectable(false)
	titleHeader := tview.NewTableCell("Name").
		SetTextColor(utils.MangaPageTitleColor).
		SetSelectable(false)
	downloadHeader := tview.NewTableCell("Download Status").
		SetTextColor(utils.MangaPageDownloadStatColor).
		SetSelectable(false)
	scanGroupHeader := tview.NewTableCell("ScanGroup").
		SetTextColor(utils.MangaPageScanGroupColor).
		SetSelectable(false)
	readMarkerHeader := tview.NewTableCell("Read Status").
		SetTextColor(utils.MangaPageReadStatColor).
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
		SetBordersColor(utils.MangaPageTableBorderColor).
		SetTitle("Chapters").
		SetTitleColor(utils.MangaPageTableTitleColor).
		SetBorder(true)

	// Add info and table to the grid. Set the focus to the chapter table.
	grid.AddItem(info, 0, 0, 5, 15, 0, 0, false).
		AddItem(table, 5, 0, 10, 15, 0, 0, true).
		AddItem(info, 0, 0, 15, 5, 0, 80, false).
		AddItem(table, 0, 5, 15, 10, 0, 80, true)

	ctx, cancel := context.WithCancel(context.Background())
	mangaPage := &MangaPage{
		Manga: manga,
		Grid:  grid,
		Info:  info,
		Table: table,
		sWrap: &utils.SelectorWrapper{
			Selection: map[int]struct{}{},
		},
		cWrap: &utils.ContextWrapper{
			Ctx:    ctx,
			Cancel: cancel,
		},
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
	status := ""
	if p.Manga.Attributes.Status != nil {
		status = *p.Manga.Attributes.Status
	}
	status = strings.Title(status)

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
	log.Println("Setting up manga page chapter table...")
	ctx, cancel := p.cWrap.ResetContext()
	// Set handlers.
	p.setHandlers(cancel)

	time.Sleep(loadDelay)
	defer cancel()

	// Show loading status so user knows it's loading.
	core.App.TView.QueueUpdateDraw(func() {
		loadingCell := tview.NewTableCell("Loading...").SetSelectable(false)
		p.Table.SetCell(1, 1, loadingCell)
	})

	// Get All chapters
	if p.cWrap.ToCancel(ctx) {
		return
	}
	chapters, err := p.getAllChapters(ctx)
	if err != nil { // If error getting chapters.
		if strings.Contains(err.Error(), contextCancelledError) {
			return
		}
		log.Println(fmt.Sprintf("Error getting manga chapters: %s", err.Error()))
		core.App.TView.QueueUpdateDraw(func() {
			modal := okModal(utils.GenericAPIErrorModalID, "Error getting manga chapters.\nCheck log for details.")
			ShowModal(utils.GenericAPIErrorModalID, modal)
		})
		return
	} else if len(chapters) == 0 { // If there are no chapters.
		core.App.TView.QueueUpdateDraw(func() {
			noResultsCell := tview.NewTableCell("No chapters!").SetSelectable(false)
			p.Table.SetCell(1, 1, noResultsCell)
		})
		return
	}

	// Get the chapter read markers.
	markers := map[string]struct{}{}
	if core.App.Client.Auth.IsLoggedIn() {
		if p.cWrap.ToCancel(ctx) {
			return
		}
		markerResponse, err := core.App.Client.Chapter.GetReadMangaChapters(p.Manga.ID)
		if err != nil {
			log.Println(fmt.Sprintf("Error getting chapter read markers: %s", err.Error()))
			core.App.TView.QueueUpdateDraw(func() {
				modal := okModal(utils.GenericAPIErrorModalID, "Error getting chapter read markers.\nCheck log for details.")
				ShowModal(utils.GenericAPIErrorModalID, modal)
			})
			return
		}
		for _, marker := range markerResponse.Data {
			markers[marker] = struct{}{}
		}
	}

	// Fill in the chapters
	for index := 0; index < len(chapters); index++ {
		if p.cWrap.ToCancel(ctx) {
			return
		}
		chapter := chapters[index]
		// Chapter Number
		chapterNumCell := tview.NewTableCell(
			fmt.Sprintf("%-6s %s", chapter.GetChapterNum(), chapter.Attributes.TranslatedLanguage)).
			SetMaxWidth(10).SetTextColor(utils.MangaPageChapNumColor).SetReference(&chapter)

		// Chapter title
		titleCell := tview.NewTableCell(fmt.Sprintf("%-30s", chapter.GetTitle())).SetMaxWidth(30).
			SetTextColor(utils.MangaPageTitleColor)

		// Chapter download status
		var downloadStatus string
		// Check for the presence of the download folder.
		if _, err = os.Stat(p.getDownloadFolder(&chapter)); err == nil {
			downloadStatus = "Y"
		}
		downloadCell := tview.NewTableCell(downloadStatus).SetTextColor(utils.MangaPageDownloadStatColor)

		// Scanlation group
		var scanGroup string
		for _, relation := range chapter.Relationships {
			if relation.Type == mangodex.ScanlationGroupRel {
				scanGroup = relation.Attributes.(*mangodex.ScanlationGroupAttributes).Name
				break
			}
		}
		scanGroupCell := tview.NewTableCell(fmt.Sprintf("%-15s", scanGroup)).SetMaxWidth(15).
			SetTextColor(utils.MangaPageScanGroupColor)

		// Read marker
		var read string
		if !core.App.Client.Auth.IsLoggedIn() {
			read = "Not logged in!"
		} else if _, ok := markers[chapter.ID]; ok {
			read = readStatus
		}
		readCell := tview.NewTableCell(read).SetTextColor(utils.MangaPageReadStatColor)

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
	}
	core.App.TView.QueueUpdateDraw(func() {
		p.Table.Select(1, 0)
		p.Table.ScrollToBeginning()
	})
}

// getAllChapters : Get All chapters for the manga.
func (p *MangaPage) getAllChapters(ctx context.Context) ([]mangodex.Chapter, error) {
	var (
		params     = p.setGetChaptersParams()
		chapters   []mangodex.Chapter
		currOffset = 0
	)
	for {
		if p.cWrap.ToCancel(ctx) {
			return []mangodex.Chapter{}, fmt.Errorf(contextCancelledError)
		}
		params.Set("offset", strconv.Itoa(currOffset))
		list, err := core.App.Client.Chapter.GetMangaChapters(p.Manga.ID, *params)
		if err != nil {
			return []mangodex.Chapter{}, err
		}
		log.Printf("Got %d of %d chapters\n", currOffset, list.Total)
		chapters = append(chapters, list.Data...)
		currOffset += chapterOffsetRange
		if currOffset >= list.Total {
			break
		}
	}
	return chapters, nil
}

// setGetChaptersParams : Helper function to set up query parameters for getting chapters.
func (p *MangaPage) setGetChaptersParams() *url.Values {
	// Set up query parameters.
	params := url.Values{}

	// Set limits
	params.Set("limit", strconv.Itoa(chapterOffsetRange))

	// Set All chapters with user's specified languages
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

	return &params
}

// markSelected : Mark a chapter as being selected by the user on the main page table.
func (p *MangaPage) markSelected(row int) {
	chapterCell := p.Table.GetCell(row, 0)
	chapterCell.SetTextColor(tcell.ColorBlack).SetBackgroundColor(utils.MangaPageHighlightColor)

	// Add to the Selection wrapper
	p.sWrap.AddSelection(row)
}

// markUnselected : Mark a chapter as being unselected by the user on the main page table.
func (p *MangaPage) markUnselected(row int) {
	chapterCell := p.Table.GetCell(row, 0)
	chapterCell.SetTextColor(utils.MangaPageChapNumColor).SetBackgroundColor(tcell.ColorBlack)

	// Remove from the Selection wrapper
	p.sWrap.RemoveSelection(row)
}

// markAll : Marks All rows as selected or unselected.
func (p *MangaPage) markAll() {
	if p.sWrap.All {
		for row := 1; row < p.Table.GetRowCount(); row++ {
			p.markUnselected(row)
		}
	} else {
		for row := 1; row < p.Table.GetRowCount(); row++ {
			p.markSelected(row)
		}
	}
	p.sWrap.All = !p.sWrap.All
}
