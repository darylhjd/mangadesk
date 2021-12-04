package pages

/*
Help Page shows the keybindings for the application.
*/

import (
	"fmt"

	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/core"
)

// ShowHelpPage : Show the help page to the user.
func ShowHelpPage(pages *tview.Pages) {
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
		fmt.Sprintf("%-15s:%15s\n\n", "Enter", "Start download") +
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

	// Create a new grid for the text view so we can align it to the center.
	grid := tview.NewGrid().SetColumns(0, 0, 0, 0).SetRows(0, 0, 0, 0).
		AddItem(help, 0, 0, 4, 4, 0, 0, false).
		AddItem(help, 1, 1, 2, 2, 45, 100, false)

	// Set up input capture for the help page.
	SetHelpPageHandlers(pages, grid)

	pages.AddPage(HelpPageID, grid, true, false)
	g.App.SetFocus(grid)
	pages.SwitchToPage(HelpPageID)
}
