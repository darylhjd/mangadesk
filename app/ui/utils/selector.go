package utils

// SelectorWrapper : A wrapper to store selections. Used by the manga page to
// keep track of selections.
type SelectorWrapper struct {
	Selection   map[int]struct{} // Keep track of which chapters have been selected by user.
	All         bool             // Keep track of whether user has selected All or not.
	VisualStart int              // Keeps track of the start of the visual selection. -1 If none.
}

// HasSelections : Checks whether there are currently selections.
func (s *SelectorWrapper) HasSelections() bool {
	return len(s.Selection) != 0
}

// HasSelection : Checks whether the current row is selected.
func (s *SelectorWrapper) HasSelection(row int) bool {
	_, ok := s.Selection[row]
	return ok
}

// CopySelection : Returns a copy of the current Selection.
func (s *SelectorWrapper) CopySelection(row int) map[int]struct{} {
	// If there are no selections currently, we add current row as a selection.
	if !s.HasSelections() {
		s.AddSelection(row)
	}
	selection := map[int]struct{}{}
	for se := range s.Selection {
		selection[se] = struct{}{}
	}

	// If there was only selection, then we treat it as a one-off transaction, and reset the current selection.
	if len(s.Selection) == 1 {
		s.Selection = map[int]struct{}{}
	}

	return selection
}

// AddSelection : Add a row to the Selection.
func (s *SelectorWrapper) AddSelection(row int) {
	s.Selection[row] = struct{}{}
}

// RemoveSelection : Remove a row from the Selection. No-op if row is not originally in Selection.
func (s *SelectorWrapper) RemoveSelection(row int) {
	delete(s.Selection, row)
}

func (s *SelectorWrapper) IsInVisualMode() bool {
	return s.VisualStart != -1
}

func (s *SelectorWrapper) StartVisualSelection(row int) {
	s.VisualStart = row
}

func (s *SelectorWrapper) StopVisualSelection() {
	s.VisualStart = -1
}
