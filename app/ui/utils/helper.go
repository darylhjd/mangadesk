package utils

import (
	"github.com/rivo/tview"
)

// NewGrid : Create a new grid with specified rows and columns.
func NewGrid(row, col []int) *tview.Grid {
	grid := tview.NewGrid().SetRows(row...).SetColumns(col...)
	return grid
}
