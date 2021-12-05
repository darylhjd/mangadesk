package ui

import "github.com/rivo/tview"

const (
	OffsetRange = 100
)

// newGrid : Create a new grid with specified rows and columns.
func newGrid(row, col []int) *tview.Grid {
	grid := tview.NewGrid().SetRows(row...).SetColumns(col...)
	return grid
}
