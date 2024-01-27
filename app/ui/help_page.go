package ui

import (
	"fmt"

	"github.com/darylhjd/mangadesk/app/core"
	"github.com/darylhjd/mangadesk/app/ui/utils"
	"github.com/rivo/tview"
)

const (
	padding = 20
)

// HelpPage : This struct contains the grid for the help page.
type HelpPage struct {
	Grid *tview.Grid
}

// ShowHelpPage : Make the app show the help page.
func ShowHelpPage() {
	helpPage := newHelpPage()

	core.App.TView.SetFocus(helpPage.Grid)
	core.App.PageHolder.AddPage(utils.HelpPageID, helpPage.Grid, true, true)
}

// newHelpPage : Creates a new help page.
func newHelpPage() *HelpPage {
	formatString := fmt.Sprintf("%%-%ds:%%%ds\n", padding, padding)
	// Set up the help text.
	helpText := "Keyboard Mappings\n" +
		"-----------------------------\n\n" +
		"Universal\n" +
		fmt.Sprintf(formatString, "Ctrl + L", "Login/Logout") +
		fmt.Sprintf(formatString, "Ctrl + K", "Keybinds/Help") +
		fmt.Sprintf(formatString, "Ctrl + S", "Search") +
		"\nManga Page\n" +
		fmt.Sprintf(formatString, "Ctrl + E", "Select mult.") +
		fmt.Sprintf(formatString, "Ctrl + A", "Toggle All") +
		fmt.Sprintf(formatString, "Ctrl + R", "Toggle Read Status") +
		fmt.Sprintf(formatString, "Ctrl + Q", "Toggle Follow Manga") +
		fmt.Sprintf(formatString, "Enter", "Start download") +
		"\nOthers\n" +
		fmt.Sprintf(formatString, "Esc", "Go back") +
		fmt.Sprintf(formatString, "Ctrl + F/B", "Next/Prev Page") +
		"\nApp Info\n" +
		core.AppVersion

	// Create TextView to show the help information.
	help := tview.NewTextView()
	// Set TextView attributes.
	help.SetText(helpText).
		SetTextAlign(tview.AlignCenter).
		SetBorderColor(utils.HelpPageBorderColor).
		SetBorder(true)

	// Create a new grid for the text view, so we can align it to the center.
	dimensions := []int{-1, -1, -1, -1, -1, -1}
	grid := utils.NewGrid(dimensions, dimensions).
		AddItem(help, 0, 0, 6, 6, 0, 0, true)

	helpPage := &HelpPage{
		Grid: grid,
	}
	helpPage.setHandlers()

	return helpPage
}
