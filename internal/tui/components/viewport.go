package components

// Viewport manages cursor position and scroll offset for a list view.
type Viewport struct {
	Cursor int
	Offset int
	Height int
	Total  int
}

// NewViewport creates a new viewport with the given height and total items.
func NewViewport(height, total int) *Viewport {
	return &Viewport{
		Height: height,
		Total:  total,
	}
}

// MoveUp moves the cursor up by one position.
func (v *Viewport) MoveUp() {
	if v.Cursor > 0 {
		v.Cursor--
		v.EnsureVisible()
	}
}

// MoveDown moves the cursor down by one position.
func (v *Viewport) MoveDown() {
	if v.Cursor < v.Total-1 {
		v.Cursor++
		v.EnsureVisible()
	}
}

// PageUp moves the cursor up by one page.
func (v *Viewport) PageUp() {
	v.Cursor = max(0, v.Cursor-v.Height)
	v.EnsureVisible()
}

// PageDown moves the cursor down by one page.
func (v *Viewport) PageDown() {
	v.Cursor = min(v.Total-1, v.Cursor+v.Height)
	v.EnsureVisible()
}

// GoToStart moves the cursor to the first item.
func (v *Viewport) GoToStart() {
	v.Cursor = 0
	v.EnsureVisible()
}

// GoToEnd moves the cursor to the last item.
func (v *Viewport) GoToEnd() {
	v.Cursor = max(0, v.Total-1)
	v.EnsureVisible()
}

// EnsureVisible adjusts the offset so the cursor is visible.
func (v *Viewport) EnsureVisible() {
	if v.Height <= 0 {
		return
	}
	if v.Cursor < v.Offset {
		v.Offset = v.Cursor
	}
	if v.Cursor >= v.Offset+v.Height {
		v.Offset = v.Cursor - v.Height + 1
	}
	v.ClampOffset()
}

// ClampOffset ensures the offset is within valid bounds.
func (v *Viewport) ClampOffset() {
	if v.Total == 0 {
		v.Offset = 0
		return
	}
	maxOffset := max(0, v.Total-v.Height)
	v.Offset = max(0, min(v.Offset, maxOffset))
}

// VisibleRange returns the start and end indices of visible items.
func (v *Viewport) VisibleRange() (start, end int) {
	start = v.Offset
	end = min(v.Offset+v.Height, v.Total)
	return start, end
}

// SetTotal updates the total item count and adjusts cursor/offset if needed.
func (v *Viewport) SetTotal(total int) {
	v.Total = total
	if v.Cursor >= total && total > 0 {
		v.Cursor = total - 1
	}
	v.ClampOffset()
	v.EnsureVisible()
}

// SetHeight updates the viewport height and adjusts offset if needed.
func (v *Viewport) SetHeight(height int) {
	v.Height = height
	v.ClampOffset()
	v.EnsureVisible()
}
