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
// Note that this does not check the current context stored in the wrapper.
func (c *ContextWrapper) toCancel(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// SelectorWrapper : A wrapper to store selections. Used by the manga page to
// keep track of selections.
type SelectorWrapper struct {
	selection map[int]struct{} // Keep track of which chapters have been selected by user.
	all       bool             // Keep track of whether user has selected all or not.
}

// hasSelections : Checks whether there are currently selections.
func (s *SelectorWrapper) hasSelections() bool {
	return len(s.selection) != 0
}

// hasSelection : Checks whether the current row is selected.
func (s *SelectorWrapper) hasSelection(row int) bool {
	_, ok := s.selection[row]
	return ok
}

// copySelection : Returns a copy of the current selection.
func (s *SelectorWrapper) copySelection() map[int]struct{} {
	selection := map[int]struct{}{}
	for se := range s.selection {
		selection[se] = struct{}{}
	}
	return selection
}

// addSelection : Add a row to the selection.
func (s *SelectorWrapper) addSelection(row int) {
	s.selection[row] = struct{}{}
}

// removeSelection : Remove a row from the selection. No-op if row is not originally in selection.
func (s *SelectorWrapper) removeSelection(row int) {
	delete(s.selection, row)
}

// newGrid : Create a new grid with specified rows and columns.
func newGrid(row, col []int) *tview.Grid {
	grid := tview.NewGrid().SetRows(row...).SetColumns(col...)
	return grid
}
