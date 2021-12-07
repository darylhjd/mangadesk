package ui

import (
	"context"
	"github.com/rivo/tview"
)

// newGrid : Create a new grid with specified rows and columns.
func newGrid(row, col []int) *tview.Grid {
	grid := tview.NewGrid().SetRows(row...).SetColumns(col...)
	return grid
}

// toCancel : Returns true if to cancel current context, false otherwise.
func toCancel(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
