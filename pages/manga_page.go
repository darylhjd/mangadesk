package pages

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

type MangaPage struct {
	InfoView     *tview.TextView
	ChapterTable *tview.Table
	Selected     *map[int]struct{} // Keep track of which chapters have been selected by user.
	SelectedAll  bool              // Keep track of if the user wants to select all.
}

// ShowMangaPage : Show the manga page.
func ShowMangaPage(pages *tview.Pages, mr *mangodex.MangaResponse) {
	// Create the base main grid.
	// 15x15 grid.
	var ga []int
	for i := 0; i < 15; i++ { // This is to create 15 grids.
		ga = append(ga, -1)
	}
	grid := tview.NewGrid().SetColumns(ga...).SetRows(ga...)
	// Set grid attributes
	grid.SetTitleColor(g.MangaPageGridTitleColor).
		SetBorderColor(g.MangaPageGridBorderColor).
		SetTitle("Manga Information").
		SetBorder(true)

	// Use a TextView for basic information of the manga.
	info := tview.NewTextView()
	// Set textview attributes
	info.SetWrap(true).SetWordWrap(true).
		SetBorderColor(g.MangaPageInfoViewBorderColor).
		SetTitleColor(g.MangaPageInfoViewTitleColor).
		SetTitle("About").
		SetBorder(true)

	// Use a table to show the chapters for the manga.
	table := tview.NewTable()
	// Set chapter headers
	numHeader := tview.NewTableCell("Chap").
		SetTextColor(g.MangaPageChapNumColor).
		SetSelectable(false)
	titleHeader := tview.NewTableCell("Name").
		SetTextColor(g.MangaPageTitleColor).
		SetSelectable(false)
	downloadHeader := tview.NewTableCell("Download Status").
		SetTextColor(g.MangaPageDownloadStatColor).
		SetSelectable(false)
	readMarkerHeader := tview.NewTableCell("Read Status").
		SetTextColor(g.MangaPageReadStatColor).
		SetSelectable(false)
	table.SetCell(0, 0, numHeader).
		SetCell(0, 1, titleHeader).
		SetCell(0, 2, downloadHeader).
		SetCell(0, 3, readMarkerHeader).
		SetFixed(1, 0)
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(g.MangaPageTableBorderColor).
		SetTitle("Chapters").
		SetTitleColor(g.MangaPageTableTitleColor).
		SetBorder(true)

	mangaPage := MangaPage{
		InfoView:     info,
		ChapterTable: table,
		Selected:     &map[int]struct{}{},
		SelectedAll:  false,
	}

	// Set input handlers for the manga page.
	// Use context to stop any goroutines that are no longer needed.
	// The page handler ESC will induce cancel.
	ctx, cancel := context.WithCancel(context.Background())
	SetMangaPageHandlers(cancel, pages, grid)

	// Set up manga info and chapter info.
	go func() {
		mangaPage.SetMangaInfo(ctx, mr)
		mangaPage.SetChapterTable(ctx, pages, mr)
	}()

	// Add info and table to the grid. Set the focus to the chapter table.
	grid.AddItem(mangaPage.InfoView, 0, 0, 5, 15, 0, 0, false).
		AddItem(mangaPage.ChapterTable, 5, 0, 10, 15, 0, 0, true).
		AddItem(mangaPage.InfoView, 0, 0, 15, 5, 0, 80, false).
		AddItem(mangaPage.ChapterTable, 0, 5, 15, 10, 0, 80, true)

	pages.AddPage(g.MangaPageID, grid, true, false)
	g.App.SetFocus(grid)
	pages.SwitchToPage(g.MangaPageID)
}

// SetMangaInfo : Populate the info TextView with required information.
// NOTE: This is run as a GOROUTINE. Drawing will require QueueUpdateDraw.
func (mp *MangaPage) SetMangaInfo(ctx context.Context, mr *mangodex.MangaResponse) {
	// Get author information
	author := "-"
CheckRelationshipLoop:
	for _, r := range mr.Relationships {
		select {
		case <-ctx.Done():
			return
		default:
			if r.Type != "author" {
				continue
			}

			if r.ID == "" {
				break CheckRelationshipLoop
			}
			if a, err := g.Dex.GetAuthor(r.ID); err != nil {
				author = "?"
			} else {
				author = a.Data.Attributes.Name
			}
			break CheckRelationshipLoop
		}
	}

	// Get status information
	status := "-"
	if mr.Data.Attributes.Status != nil {
		status = strings.Title(*mr.Data.Attributes.Status)
	}

	// Set up information text.
	infoText := fmt.Sprintf("Title: %s\n\nAuthor: %s\nStatus: %s\n\nDescription:\n%s",
		mr.Data.Attributes.Title["en"], author, status,
		strings.SplitN(tview.Escape(mr.Data.Attributes.Description["en"]), "\n", 2)[0])

	g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
		mp.InfoView.SetText(infoText)
	})
}

// SetChapterTable : Populate the manga page chapter table.
// NOTE: This is run as a GOROUTINE. Drawing will require QueueUpdateDraw.
func (mp *MangaPage) SetChapterTable(ctx context.Context, pages *tview.Pages, mr *mangodex.MangaResponse) {
	// Get all chapters for this manga. No pages.
	chapters, ok := mp.GetAllChapters(ctx, pages, mr.Data.ID)
	if !ok {
		return
	} else if len(*chapters) == 0 { // If no chapters
		noResCell := tview.NewTableCell("No chapters.").SetSelectable(false)
		g.App.QueueUpdateDraw(func() {
			mp.ChapterTable.SetCell(1, 1, noResCell)
		})
		return
	}

	// Set input handlers for the table
	mp.ChapterTable.SetSelectedFunc(func(row, column int) { // When user presses ENTER to confirm selected.
		// We add the current selection if the there are no selected rows currently.
		if len(*mp.Selected) == 0 {
			(*mp.Selected)[row] = struct{}{}
		}
		// Show modal to confirm download.
		ShowModal(pages, g.DownloadChaptersModalID, "Download selection(s)?", []string{"Yes", "No"},
			func(i int, label string) {
				if label == "Yes" {
					// If user confirms to download, then we download the chapters.
					downloadChapters(pages, mp, mr, chapters)
				}
				pages.RemovePage(g.DownloadChaptersModalID)
			})
	})
	SetMangaPageTableHandlers(mp, len(*chapters)) // For custom input handlers.

	// Add each chapter info to the table.
	for i, cr := range *chapters {
		select {
		case <-ctx.Done():
			return
		default:
			// Chapter number cell.
			c := "-"
			if cr.Data.Attributes.Chapter != nil {
				c = *cr.Data.Attributes.Chapter
			}
			cCell := tview.NewTableCell(fmt.Sprintf("%-6s", c)).SetMaxWidth(6).
				SetTextColor(g.MangaPageChapNumColor)

			// Chapter title cell.
			tCell := tview.NewTableCell(fmt.Sprintf("%-40s", cr.Data.Attributes.Title)).SetMaxWidth(40).
				SetTextColor(g.MangaPageTitleColor)

			// Chapter download status cell.
			// Get chapter folder name.
			chapter := generateChapterFolderName(&cr)
			chapFolder := filepath.Join(g.Conf.DownloadDir, mr.Data.Attributes.Title["en"], chapter)

			// Check whether the folder for this chapter exists. If it does, then it is downloaded.
			stat := ""
			if _, err := os.Stat(chapFolder); err == nil {
				stat = "Yup!"
			}
			dCell := tview.NewTableCell(stat).SetTextColor(g.MangaPageReadStatColor)

			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				mp.ChapterTable.SetCell(i+1, 0, cCell)
				mp.ChapterTable.SetCell(i+1, 1, tCell)
				mp.ChapterTable.SetCell(i+1, 2, dCell)
			})
		}
	}

	// Set read markers for chapters.
	mp.SetChapterReadMarkers(ctx, mr.Data.ID, chapters)
}

// GetAllChapters : Get all chapters for a manga.
// NOTE: This is run in a GOROUTINE. Drawing will require QueueUpdateDraw.
func (mp *MangaPage) GetAllChapters(ctx context.Context, pages *tview.Pages, mangaID string) (*[]mangodex.ChapterResponse, bool) {
	// Set up query parameters to get chapters.
	params := url.Values{}
	params.Set("limit", "500")
	for _, lang := range g.Conf.Languages { // Add user's languages
		params.Add("translatedLanguage[]", lang)
	}
	params.Set("order[chapter]", "desc") // Show latest chapters first

	var (
		chapters []mangodex.ChapterResponse
		offset   = 0
	)
GetAllChapterLoop:
	for {
		select {
		case <-ctx.Done(): // If user already exited the manga page.
			return nil, false
		default:
			params.Set("offset", strconv.Itoa(offset))
			chapterList, err := g.Dex.MangaFeed(mangaID, params)
			if err != nil {
				// If error getting chapters for the manga, we tell the user so through a modal.
				g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
					OKModal(pages, g.GenericAPIErrorModalID, "Error getting manga feed.")
				})
				return nil, false // We end immediately. No need to continue.
			}
			chapters = append(chapters, chapterList.Results...)
			// Check if there are still more chapters to load.
			offset += 500
			if offset >= chapterList.Total {
				break GetAllChapterLoop
			}
		}
	}
	return &chapters, true
}

// SetChapterReadMarkers : Set read markers for each chapter for a manga.
// NOTE: This is run in a GOROUTINE. Drawing will require QueueUpdateDraw.
func (mp *MangaPage) SetChapterReadMarkers(ctx context.Context, mangaID string, chapters *[]mangodex.ChapterResponse) {
	// Check for manga read markers.
	if !g.Dex.IsLoggedIn() { // If user is not logged in.
		// We inform user to log in to track read status.
		// Split the message into 2 rows.
		rSCell1 := tview.NewTableCell("Log in to").SetTextColor(g.MangaPageReadStatColor)
		rSCell2 := tview.NewTableCell("see read status!").SetTextColor(g.MangaPageReadStatColor)

		g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
			mp.ChapterTable.SetCell(1, 3, rSCell1)
			mp.ChapterTable.SetCell(2, 3, rSCell2)
		})
		return // We return immediately. No need to continue.
	}

	// Use a map to store the read chapter IDs to avoid iterating through every turn.
	read := map[string]struct{}{}
	select {
	case <-ctx.Done():
		return
	default:
		chapReadMarkerResp, err := g.Dex.MangaReadMarkers(mangaID)
		if err != nil { // If error getting read markers, just put a error message on the column.
			readStatus := "API Error!"
			g.App.QueueUpdateDraw(func() {
				rSCell := tview.NewTableCell(readStatus).SetTextColor(g.MangaPageReadStatColor)
				mp.ChapterTable.SetCell(1, 3, rSCell)
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
			if _, ok := read[cr.Data.ID]; ok { // If chapter ID is in map of read markers.
				readStatus = "R"
			}
			rSCell := tview.NewTableCell(readStatus).SetTextColor(g.MangaPageReadStatColor)
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				mp.ChapterTable.SetCell(i+1, 3, rSCell)
			})
		}
	}
}

// MarkChapterSelected : Mark a chapter as being selected by the user on the main page table.
func markChapterSelected(table *tview.Table, row int) {
	chapterCell := table.GetCell(row, 0)
	chapterCell.SetTextColor(tcell.ColorBlack).SetBackgroundColor(g.MangaPageHighlightColor)
}

// MarkChapterUnselected : Mark a chapter as being unselected by the user on the main page table.
func markChapterUnselected(table *tview.Table, row int) {
	chapterCell := table.GetCell(row, 0)
	chapterCell.SetTextColor(g.MangaPageChapNumColor).SetBackgroundColor(tcell.ColorBlack)
}
