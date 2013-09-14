// sheet-movement
package sheet

func (s *Sheet) MoveUp() {
	sel := s.SelectedCell
	row, colIdx := sel.RowCol()
	row--
	if row < 0 {
		row = 0
	}
	if int(row) < s.startRow {
		s.startRow = int(row)
	}
	s.SelectedCell = getAddress(row, colIdx)
	s.display()
}

func (s *Sheet) MoveDown() {
	sel := s.SelectedCell
	row, colIdx := sel.RowCol()
	row++
	s.SelectedCell = getAddress(row, colIdx)
	if (int(row)-s.startRow)+1 > s.displayRows {
		s.startRow++
	}
	s.display()
}

func (s *Sheet) MoveRight() {
	sel := s.SelectedCell
	row, colIdx := sel.RowCol()
	colIdx++
	if (colIdx-s.startCol)+1 > s.displayCols {
		s.startCol++
	}
	s.SelectedCell = getAddress(row, colIdx)
	s.display()
}

func (s *Sheet) MoveLeft() {
	sel := s.SelectedCell
	row, colIdx := sel.RowCol()
	colIdx--
	if colIdx < 0 {
		colIdx = 0
	}
	if colIdx < s.startCol {
		s.startCol = colIdx
	}
	s.SelectedCell = getAddress(row, colIdx)
	s.display()
}
