package pages

import (
	"context"
	"fmt"
	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/darylhjd/mangadesk/core"
)

// MangaPage : This struct contains the required primitives for the manga page.
type MangaPage struct {
	Manga       *mangodex.Manga
	Grid        *tview.Grid
	Info        *tview.TextView
	Table       *tview.Table
	Selected    map[int]struct{} // Keep track of which chapters have been selected by user.
	SelectedAll bool             // Keep track of if the user wants to select all.
}

// NewMangaPage : Creates a new manga page.
func NewMangaPage(manga *mangodex.Manga) *MangaPage {
	var dimensions []int
	for i := 0; i < 15; i++ {
		dimensions = append(dimensions, -1)
	}
	grid := NewGrid(dimensions, dimensions)
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
	readMarkerHeader := tview.NewTableCell("Read Status").
		SetTextColor(MangaPageReadStatColor).
		SetSelectable(false)
	scanGroupHeader := tview.NewTableCell("ScanGroup").
		SetTextColor(MangaPageScanGroupColor).
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
		Manga: manga,
		Grid:  grid,
		Info:  info,
		Table: table,
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
		}
	}

	// Status
	status := *p.Manga.Attributes.Status

	// Description
	desc := tview.Escape(p.Manga.GetDescription("en"))

	// Set up information text.
	infoText := fmt.Sprintf("Title: %s\n\nAuthor: %s\nStatus: %s\n\nDescription:\n%s",
		title, author, status, desc)

	core.App.ViewApp.QueueUpdateDraw(func() {
		p.Info.SetText(infoText)
	})
}

// setChapterTable : Fill up the chapter table.
func (p *MangaPage) setChapterTable() {
	// Show loading status so user knows it's loading.
	core.App.ViewApp.QueueUpdateDraw(func() {
		loadingCell := tview.NewTableCell("Loading...").SetSelectable(false)
		p.Table.SetCell(1, 1, loadingCell)
	})

	// Get all chapters
	chapters, err := p.getAllChapters()
	if err != nil {
		log.Println(fmt.Sprintf("Error getting manga chapters: %s", err.Error()))
		modal := OKModal(GenericAPIErrorModalID, "Error getting manga chapters.\nCheck log for details.")
		core.App.ShowModal(GenericAPIErrorModalID, modal)
		return
	}

	for index, chapter := range chapters {
		chapterNumCell := tview.NewTableCell(
			fmt.Sprintf("%-6s %s", chapter.GetChapterNum(), chapter.Attributes.TranslatedLanguage)).
			SetMaxWidth(10).SetTextColor(MangaPageChapNumColor)

		titleCell := tview.NewTableCell(fmt.Sprintf("%-30s", chapter.GetTitle())).SetMaxWidth(30).
			SetTextColor(MangaPageTitleColor)

		var downloadStatus string
		// Check for the presence of the download folder.
	}
}

// SetChapterTable : Populate the manga page chapter table.
// NOTE: This is run as a GOROUTINE. Drawing will require QueueUpdateDraw.
func (mp *MangaPage) SetChapterTable(ctx context.Context, pages *tview.Pages, m *mangodex.Manga) {
	// Add each chapter info to the table.
	for i, c := range *chapters {
		select {
		case <-ctx.Done():
			return
		default:
			// Chapter download status cell.
			// Get the manga and chapter folder name.
			mangaName, chapter := generateChapterFolderNames(m, &c)
			chapFolder := filepath.Join(core.Conf.DownloadDir, mangaName, chapter)
			// Check whether the folder for this chapter exists. If it does, then it is downloaded.
			stat := ""
			if _, err := os.Stat(chapFolder); err == nil {
				stat = "Y"
			} else if _, err = os.Stat(fmt.Sprintf("%s.%s", chapFolder, core.Conf.ZipType)); err == nil {
				stat = "Y"
			}
			dCell := tview.NewTableCell(stat).SetTextColor(MangaPageDownloadStatColor)

			core.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				mp.ChapterTable.SetCell(i+1, 0, cCell).
					SetCell(i+1, 1, tCell).
					SetCell(i+1, 2, dCell)
			})
		}
	}

	// Set scanlation groups for chapters.
	mp.SetChapterScanGroup(ctx, chapters)
	// Set read markers for chapters.
	mp.SetChapterReadMarkers(ctx, m.ID, chapters)
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

// SetChapterReadMarkers : Set read markers for each chapter for a manga.
// NOTE: This is run in a GOROUTINE. Drawing will require QueueUpdateDraw.
func (mp *MangaPage) SetChapterReadMarkers(ctx context.Context, mangaID string, chapters *[]mangodex.Chapter) {
	// Check for manga read markers.
	if !core.DexClient.IsLoggedIn() { // If user is not logged in.
		// We inform user to log in to track read status.
		// Split the message into 2 rows.
		rSCell := tview.NewTableCell("Not logged in!").SetTextColor(MangaPageReadStatColor).SetSelectable(false)

		core.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
			mp.ChapterTable.SetCell(1, 4, rSCell)
		})
		return // We return immediately. No need to continue.
	}

	// Use a map to store the read chapter IDs to avoid iterating through every turn.
	read := map[string]struct{}{}
	select {
	case <-ctx.Done():
		return
	default:
		chapReadMarkerResp, err := core.DexClient.MangaReadMarkers(mangaID)
		if err != nil { // If error getting read markers, just put a error message on the column.
			readStatus := "API Error!"
			core.App.QueueUpdateDraw(func() {
				rSCell := tview.NewTableCell(readStatus).SetTextColor(MangaPageReadStatColor)
				mp.ChapterTable.SetCell(1, 4, rSCell)
			})
			return // We return immediately. No need to continue.
		}
		for _, chapID := range chapReadMarkerResp.Data {
			read[chapID] = struct{}{}
		}
	}

	// For every chapter
	for i, cr := range *chapters {
		select {
		case <-ctx.Done():
			return
		default:
			readStatus := ""
			if _, ok := read[cr.ID]; ok { // If chapter ID is in map of read markers.
				readStatus = "R"
			}
			rSCell := tview.NewTableCell(readStatus).SetTextColor(MangaPageReadStatColor)
			core.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				mp.ChapterTable.SetCell(i+1, 4, rSCell)
			})
		}
	}
}

func (mp *MangaPage) SetChapterScanGroup(ctx context.Context, chapters *[]mangodex.Chapter) {
	// Keep track of already checked groups to avoid calling the API too many times.
	groups := map[string]string{}
	for i, c := range *chapters {
		select {
		case <-ctx.Done():
			return
		default:
			// Find the scanlation_group relationship
			group := ""
			for _, r := range c.Relationships {
				if r.Type == "scanlation_group" {
					groupId := r.ID
					if name, ok := groups[groupId]; ok {
						group = name
					} else {
						sgr, err := core.DexClient.ViewScanGroup(groupId)
						if err != nil {
							group = "API Error!"
							core.App.QueueUpdateDraw(func() {
								sgCell := tview.NewTableCell(group).SetTextColor(MangaPageScanGroupColor)
								mp.ChapterTable.SetCell(1, 3, sgCell)
							})
							return
						}
						group = sgr.Data.Attributes.Name
						// Add to cache
						groups[groupId] = group
					}
					// Get group name using ID
					sgCell := tview.NewTableCell(fmt.Sprintf("%-15s", group)).SetMaxWidth(15).
						SetTextColor(MangaPageScanGroupColor)
					core.App.QueueUpdateDraw(func() {
						mp.ChapterTable.SetCell(i+1, 3, sgCell)
					})
					break
				}
			}
		}
	}
}

// MarkChapterSelected : Mark a chapter as being selected by the user on the main page table.
func markChapterSelected(table *tview.Table, row int) {
	chapterCell := table.GetCell(row, 0)
	chapterCell.SetTextColor(tcell.ColorBlack).SetBackgroundColor(MangaPageHighlightColor)
}

// MarkChapterUnselected : Mark a chapter as being unselected by the user on the main page table.
func markChapterUnselected(table *tview.Table, row int) {
	chapterCell := table.GetCell(row, 0)
	chapterCell.SetTextColor(MangaPageChapNumColor).SetBackgroundColor(tcell.ColorBlack)
}
