package pages

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

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
	grid.SetTitleColor(tcell.ColorOrange).
		SetBorderColor(tcell.ColorLightGrey).
		SetTitle("Manga Information").
		SetBorder(true)

	// Set input handlers for the manga page.
	SetMangaPageHandlers(pages, grid)

	// Use a TextView for basic information of the manga.
	info := tview.NewTextView()
	// Set manga info on the TextView.
	go func() {
		setMangaInfo(info, mr)
	}()
	// Set textview attributes
	info.SetWrap(true).SetWordWrap(true).
		SetBorderColor(tcell.ColorLightGrey).
		SetTitleColor(tcell.ColorLightSkyBlue).
		SetTitle("About").
		SetBorder(true)

	// Use a table to show the chapters for the manga.
	table := tview.NewTable()
	// Set chapter headers
	numHeader := tview.NewTableCell("Chap").
		SetTextColor(tcell.ColorLightYellow).
		SetSelectable(false)
	titleHeader := tview.NewTableCell("Name").
		SetTextColor(tcell.ColorLightSkyBlue).
		SetSelectable(false)
	downloadHeader := tview.NewTableCell("Download Status").
		SetTextColor(tcell.ColorPowderBlue).
		SetSelectable(false)
	readMarkerHeader := tview.NewTableCell("Read Status").
		SetTextColor(tcell.ColorOrange).
		SetSelectable(false)
	table.SetCell(0, 0, numHeader).
		SetCell(0, 1, titleHeader).
		SetCell(0, 2, downloadHeader).
		SetCell(0, 3, readMarkerHeader).
		SetFixed(1, 0)
	// Set up chapter info on the table.
	go func() {
		setMangaChaptersTable(pages, table, mr)
	}()
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(tcell.ColorGrey).
		SetTitle("Chapters").
		SetTitleColor(tcell.ColorLightSkyBlue).
		SetBorder(true)

	// Add info and table to the grid. Set the focus to the chapter table.
	grid.AddItem(info, 0, 0, 5, 15, 0, 0, false).
		AddItem(table, 5, 0, 10, 15, 0, 0, true).
		AddItem(info, 0, 0, 15, 5, 0, 80, false).
		AddItem(table, 0, 5, 15, 10, 0, 80, true)

	pages.AddPage(g.MangaPageID, grid, true, false)
	g.App.SetFocus(grid)
	pages.SwitchToPage(g.MangaPageID)
}

// setMangaInfo : Populate the info TextView with required information.
// NOTE: This is run as a GOROUTINE. Drawing will require QueueUpdateDraw.
func setMangaInfo(info *tview.TextView, mr *mangodex.MangaResponse) {
	// Get author information
	authorId := ""
	for _, r := range mr.Relationships {
		if r.Type == "author" {
			authorId = r.ID
			break
		}
	}
	author := "-"
	if authorId != "" {
		a, err := g.Dex.GetAuthor(authorId)
		if err != nil {
			author = "?"
		} else {
			author = a.Data.Attributes.Name
		}
	}

	// Get status information
	status := "-"
	if mr.Data.Attributes.Status != nil {
		status = strings.Title(*mr.Data.Attributes.Status)
	}

	infoText := fmt.Sprintf("Title: %s\n\nAuthor: %s\nStatus: %s\n\nDescription:\n%s",
		mr.Data.Attributes.Title["en"], author, status,
		strings.SplitN(tview.Escape(mr.Data.Attributes.Description["en"]), "\n", 2)[0])

	g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
		info.SetText(infoText)
	})
}

// setMangaChaptersTable : Populate the manga page chapter table.
// NOTE: This is run as a GOROUTINE. Drawing will require QueueUpdateDraw.
func setMangaChaptersTable(pages *tview.Pages, table *tview.Table, mr *mangodex.MangaResponse) {
	// Get chapter feed for this manga.
	// Set up query parameters to get chapters.
	params := url.Values{}
	params.Set("limit", "500")
	params.Set("translatedLanguage[]", "en")
	params.Set("order[chapter]", "desc")
	cl, err := g.Dex.MangaFeed(mr.Data.ID, params)
	if err != nil {
		// If error getting chapters for the manga, we tell the user so through a modal.
		g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
			ShowModal(pages, g.GenericAPIErrorModalID, "Error getting manga feed", []string{"OK"},
				func(i int, label string) {
					pages.RemovePage(g.GenericAPIErrorModalID)
				})
		})
		return // We end immediately. No need to continue.
	}

	// Set input handlers for the table
	selected := map[int]struct{}{}                // We use this map to keep track of which chapters have been selected by user.
	table.SetSelectedFunc(func(row, column int) { // When user presses ENTER to confirm selected.
		confirmChapterDownloads(pages, table, &selected, row, mr, cl)
	})
	SetMangaPageTableHandlers(table, &selected) // For custom input handlers.

	// Add each chapter info to the table.
	for i, cr := range cl.Results {
		g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
			// Create chapter cell and put chapter number.
			c := "-"
			if cr.Data.Attributes.Chapter != nil {
				c = *cr.Data.Attributes.Chapter
			}
			cCell := tview.NewTableCell(fmt.Sprintf("%-5s", c)).SetMaxWidth(5).
				SetTextColor(tcell.ColorLightYellow)

			// Create title cell and put title.
			tCell := tview.NewTableCell(fmt.Sprintf("%-40s", cr.Data.Attributes.Title)).
				SetTextColor(tcell.ColorLightSkyBlue).SetMaxWidth(40)

			// Create the downloaded status cell and put the status inside.
			// Get chapter folder name.
			chapter := generateChapterFolderName(&cr)
			chapFolder := filepath.Join(g.Conf.DownloadDir, mr.Data.Attributes.Title["en"], chapter)
			stat := ""
			// Check whether the folder for this chapter exists. If it does, then it is downloaded.
			if _, err = os.Stat(chapFolder); err == nil {
				stat = "Yup!"
			}
			dCell := tview.NewTableCell(stat).SetTextColor(tcell.ColorPowderBlue)

			table.SetCell(i+1, 0, cCell)
			table.SetCell(i+1, 1, tCell)
			table.SetCell(i+1, 2, dCell)
		})
	}

	// Check for manga read markers.
	if !g.Dex.IsLoggedIn() { // If user is not logged in.
		// We inform user to log in to track read status.
		// Split the message into 2 rows.
		rSCell1 := tview.NewTableCell("Log in to").SetTextColor(tcell.ColorOrange)
		rSCell2 := tview.NewTableCell("see read status!").SetTextColor(tcell.ColorOrange)

		g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
			table.SetCell(1, 3, rSCell1)
			table.SetCell(2, 3, rSCell2)
		})
		return // We return immediately. No need to continue.
	}
	readStatus := ""
	chapReadMarkerResp, err := g.Dex.MangaReadMarkers(mr.Data.ID)
	if err != nil { // If error getting read markers, just put a error message on the column.
		readStatus = "API Error!"
		g.App.QueueUpdateDraw(func() {
			rSCell := tview.NewTableCell(readStatus).SetTextColor(tcell.ColorOrange)
			table.SetCell(1, 3, rSCell)
		})
		return // We return immediately. No need to continue.
	}
	// If no error, we can go ahead and check each chapter.
	// Use a map to store the read chapter IDs to avoid iterating through every turn.
	read := map[string]struct{}{}
	for _, chapID := range chapReadMarkerResp.Data {
		read[chapID] = struct{}{}
	}
	// For every chapter
	for i, cr := range cl.Results {
		if _, ok := read[cr.Data.ID]; ok { // If chapter ID is in map of read markers.
			readStatus = "R"
		}
		rSCell := tview.NewTableCell(readStatus).SetTextColor(tcell.ColorOrange)
		g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
			table.SetCell(i+1, 3, rSCell)
		})
	}
}

// MarkChapterSelected : Mark a chapter as being selected by the user on the main page table.
func markChapterSelected(table *tview.Table, row int) {
	chapterCell := table.GetCell(row, 0)
	chapterCell.SetBackgroundColor(tcell.ColorLimeGreen).SetTextColor(tcell.ColorBlack)

	titleCell := table.GetCell(row, 1)
	titleCell.SetBackgroundColor(tcell.ColorLimeGreen).SetTextColor(tcell.ColorBlack)
}

// MarkChapterUnselected : Mark a chapter as being unselected by the user on the main page table.
func markChapterUnselected(table *tview.Table, row int) {
	cCell := table.GetCell(row, 0)
	cCell.SetTextColor(tcell.ColorLightYellow).SetBackgroundColor(tcell.ColorBlack)

	tCell := table.GetCell(row, 1)
	tCell.SetTextColor(tcell.ColorLightSkyBlue).SetBackgroundColor(tcell.ColorBlack)
}
