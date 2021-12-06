package ui

import (
	"fmt"
	"github.com/darylhjd/mangadesk/app/core"
	"github.com/rivo/tview"
)

// HelpPage : This struct contains the grid for the help page.
type HelpPage struct {
	Grid *tview.Grid
}

// ShowHelpPage : Make the app show the help page.
func ShowHelpPage() {
	helpPage := newHelpPage()

	core.App.TView.SetFocus(helpPage.Grid)
	core.App.PageHolder.AddPage(HelpPageID, helpPage.Grid, true, true)
}

// newHelpPage : Creates a new help page.
func newHelpPage() *HelpPage {
	// Set up the help text.
	helpText := "Keyboard Mappings\n" +
		"-----------------------------\n\n" +
		"Universal\n" +
		fmt.Sprintf("%-15s:%15s\n", "Ctrl + L", "Login/Logout") +
		fmt.Sprintf("%-15s:%15s\n", "Ctrl + K", "Keybinds/Help") +
		fmt.Sprintf("%-15s:%15s\n\n", "Ctrl + S", "Search") +
		"Manga Page\n" +
		fmt.Sprintf("%-15s:%15s\n", "Ctrl + E", "Select mult.") +
		fmt.Sprintf("%-15s:%15s\n", "Ctrl + A", "Toggle All") +
		fmt.Sprintf("%-15s:%15s\n\n", "Enter", "Initialise download") +
		"Others\n" +
		fmt.Sprintf("%-15s:%15s\n", "Esc", "Go back") +
		fmt.Sprintf("%-15s:%15s\n\n", "Ctrl + F/B", "Next/Prev Page")

	// Create TextView to show the help information.
	help := tview.NewTextView()
	// Set TextView attributes.
	help.SetText(helpText).
		SetTextAlign(tview.AlignCenter).
		SetBorderColor(HelpPageBorderColor).
		SetBorder(true)

	// Create a new grid for the text view, so we can align it to the center.
	grid := tview.NewGrid().SetColumns(0, 0, 0, 0).SetRows(0, 0, 0, 0).
		AddItem(help, 0, 0, 4, 4, 0, 0, false).
		AddItem(help, 1, 1, 2, 2, 45, 100, false)

	helpPage := &HelpPage{
		Grid: grid,
	}
	helpPage.setHandlers()

	return helpPage
}
