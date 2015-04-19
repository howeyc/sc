package sheet

// Move selected cell up on row.
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
	s.SelectedCell = NewAddress(row, colIdx)
	s.display()
}

// Move selected cell down one row.
func (s *Sheet) MoveDown() {
	sel := s.SelectedCell
	row, colIdx := sel.RowCol()
	row++
	s.SelectedCell = NewAddress(row, colIdx)
	if (int(row)-s.startRow)+1 > s.displayRows {
		s.startRow++
	}
	s.display()
}

// Move selected cell right on column.
func (s *Sheet) MoveRight() {
	sel := s.SelectedCell
	row, colIdx := sel.RowCol()
	colIdx++
	if (colIdx-s.startCol)+1 > s.displayCols {
		s.startCol++
	}
	s.SelectedCell = NewAddress(row, colIdx)
	s.display()
}

// Move selected cell left one column.
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
	s.SelectedCell = NewAddress(row, colIdx)
	s.display()
}

// Move to a given address
func (s *Sheet) GoTo(adrs Address) {
	curRow, curCol := s.SelectedCell.RowCol()
	dstRow, dstCol := adrs.RowCol()

	// Move row
	for rowIdx := curRow; rowIdx < dstRow; rowIdx++ {
		s.MoveDown()
	}
	for rowIdx := curRow; rowIdx > dstRow; rowIdx-- {
		s.MoveUp()
	}

	// Move col
	for colIdx := curCol; colIdx < dstCol; colIdx++ {
		s.MoveRight()
	}
	for colIdx := curCol; colIdx > dstCol; colIdx-- {
		s.MoveLeft()
	}
}
