package pages

/*
Manga Page shows information including chapters for a particular manga.
*/

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
func ShowMangaPage(pages *tview.Pages, m *mangodex.Manga) {
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
	scanGroupHeader := tview.NewTableCell("ScanGroup").
		SetTextColor(g.MangaPageScanGroupColor).
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
		SetBordersColor(g.MangaPageTableBorderColor).
		SetTitle("Chapters").
		SetTitleColor(g.MangaPageTableTitleColor).
		SetBorder(true)

	// Create the MangaPage.
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
		defer cancel() // Hit off the cancel function if it has not yet been cancelled.
		mangaPage.SetMangaInfo(ctx, m)
		mangaPage.SetChapterTable(ctx, pages, m)
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
func (mp *MangaPage) SetMangaInfo(ctx context.Context, m *mangodex.Manga) {
	// Get author information
	author := "-"
CheckRelationshipLoop:
	for _, r := range m.Relationships {
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
	if m.Attributes.Status != nil {
		status = strings.Title(*m.Attributes.Status)
	}

	// Set up information text.
	infoText := fmt.Sprintf("Title: %s\n\nAuthor: %s\nStatus: %s\n\nDescription:\n%s",
		m.Attributes.Title["en"], author, status, tview.Escape(m.Attributes.Description["en"]))

	g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
		mp.InfoView.SetText(infoText)
	})
}

// SetChapterTable : Populate the manga page chapter table.
// NOTE: This is run as a GOROUTINE. Drawing will require QueueUpdateDraw.
func (mp *MangaPage) SetChapterTable(ctx context.Context, pages *tview.Pages, m *mangodex.Manga) {
	// Show loading status so user knows it's loading.
	g.App.QueueUpdateDraw(func() {
		loadingCell := tview.NewTableCell("Loading...").SetSelectable(false)
		mp.ChapterTable.SetCell(1, 1, loadingCell)
	})

	// Get all chapters for this manga. No pages.
	chapters, ok := mp.GetAllChapters(ctx, pages, m)
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
	SetMangaPageTableHandlers(pages, mp, m, chapters) // For custom input handlers.

	// Add each chapter info to the table.
	for i, c := range *chapters {
		select {
		case <-ctx.Done():
			return
		default:
			// Chapter number cell. Note that this column also contains the translated language.
			cNum := "-"
			if c.Attributes.Chapter != nil {
				cNum = *c.Attributes.Chapter
			}
			cCell := tview.NewTableCell(fmt.Sprintf("%-6s %s", cNum, c.Attributes.TranslatedLanguage)).
				SetMaxWidth(10).SetTextColor(g.MangaPageChapNumColor)

			// Chapter title cell.
			tCell := tview.NewTableCell(fmt.Sprintf("%-30s", c.Attributes.Title)).SetMaxWidth(30).
				SetTextColor(g.MangaPageTitleColor)

			// Chapter download status cell.
			// Get the manga and chapter folder name.
			mangaName, chapter := generateChapterFolderNames(m, &c)
			chapFolder := filepath.Join(g.Conf.DownloadDir, mangaName, chapter)
			// Check whether the folder for this chapter exists. If it does, then it is downloaded.
			stat := ""
			if _, err := os.Stat(chapFolder); err == nil {
				stat = "Y"
			} else if _, err = os.Stat(fmt.Sprintf("%s.%s", chapFolder, g.Conf.ZipType)); err == nil {
				stat = "Y"
			}
			dCell := tview.NewTableCell(stat).SetTextColor(g.MangaPageDownloadStatColor)

			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
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

// GetAllChapters : Get all chapters for a manga.
// NOTE: This is run in a GOROUTINE. Drawing will require QueueUpdateDraw.
func (mp *MangaPage) GetAllChapters(ctx context.Context, pages *tview.Pages, m *mangodex.Manga) (*[]mangodex.Chapter, bool) {
	// Set up query parameters to get chapters.
	params := url.Values{}
	params.Set("limit", "500")
	for _, lang := range g.Conf.Languages { // Add user's languages
		params.Add("translatedLanguage[]", lang)
	}
	params.Set("order[chapter]", "desc") // Show latest chapters first
	// If manga is pornographic, then also load pornographic chapters.
	if m.Attributes.ContentRating != nil && *m.Attributes.ContentRating == "pornographic" {
		ratings := []string{"safe", "suggestive", "erotica", "pornographic"}
		for _, rating := range ratings {
			params.Add("contentRating[]", rating)
		}
	}

	var (
		chapters []mangodex.Chapter
		offset   = 0
	)
GetAllChapterLoop:
	for {
		select {
		case <-ctx.Done(): // If user already exited the manga page.
			return nil, false
		default:
			params.Set("offset", strconv.Itoa(offset))
			chapterList, err := g.Dex.MangaFeed(m.ID, params)
			if err != nil {
				// If error getting chapters for the manga, we tell the user so through a modal.
				g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
					OKModal(pages, g.GenericAPIErrorModalID, "Error getting manga feed.")
				})
				return nil, false // We end immediately. No need to continue.
			}
			chapters = append(chapters, chapterList.Data...)
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
func (mp *MangaPage) SetChapterReadMarkers(ctx context.Context, mangaID string, chapters *[]mangodex.Chapter) {
	// Check for manga read markers.
	if !g.Dex.IsLoggedIn() { // If user is not logged in.
		// We inform user to log in to track read status.
		// Split the message into 2 rows.
		rSCell := tview.NewTableCell("Not logged in!").SetTextColor(g.MangaPageReadStatColor).SetSelectable(false)

		g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
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
		chapReadMarkerResp, err := g.Dex.MangaReadMarkers(mangaID)
		if err != nil { // If error getting read markers, just put a error message on the column.
			readStatus := "API Error!"
			g.App.QueueUpdateDraw(func() {
				rSCell := tview.NewTableCell(readStatus).SetTextColor(g.MangaPageReadStatColor)
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
			rSCell := tview.NewTableCell(readStatus).SetTextColor(g.MangaPageReadStatColor)
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
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
						sgr, err := g.Dex.ViewScanGroup(groupId)
						if err != nil {
							group = "API Error!"
							g.App.QueueUpdateDraw(func() {
								sgCell := tview.NewTableCell(group).SetTextColor(g.MangaPageScanGroupColor)
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
						SetTextColor(g.MangaPageScanGroupColor)
					g.App.QueueUpdateDraw(func() {
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
	chapterCell.SetTextColor(tcell.ColorBlack).SetBackgroundColor(g.MangaPageHighlightColor)
}

// MarkChapterUnselected : Mark a chapter as being unselected by the user on the main page table.
func markChapterUnselected(table *tview.Table, row int) {
	chapterCell := table.GetCell(row, 0)
	chapterCell.SetTextColor(g.MangaPageChapNumColor).SetBackgroundColor(tcell.ColorBlack)
}
