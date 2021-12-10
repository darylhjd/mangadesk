package ui

import (
	"context"
	"github.com/rivo/tview"
)

// ContextWrapper : A wrapper around a context and its corresponding cancel function.
type ContextWrapper struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// resetContext : Sets the new context, and return the old context,
// which the caller can use to cancel the previous context.
func (c *ContextWrapper) resetContext() (context.Context, context.CancelFunc) {
	ctx, cancel := c.ctx, c.cancel
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return ctx, cancel
}

// toCancel : Helper function to check if the previous context should be cancelled.
func (c *ContextWrapper) toCancel(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// newGrid : Create a new grid with specified rows and columns.
func newGrid(row, col []int) *tview.Grid {
	grid := tview.NewGrid().SetRows(row...).SetColumns(col...)
	return grid
}
