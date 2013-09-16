// sheet-commands
package sheet

// Yanks the row of the selected cell.
func (s *Sheet) YankRow() {
	row := s.SelectedCell.Row()

	s.clipboardRangeStart = NewAddress(row, 0)
	s.clipboardRangeEnd = NewAddress(row, s.maxColumnForRow[row])
}

// Yanks the column of the selected cell.
func (s *Sheet) YankColumn() {
	column := s.SelectedCell.Column()

	s.clipboardRangeStart = NewAddress(0, column)
	s.clipboardRangeEnd = NewAddress(s.maxRowForColumn[column], column)
}

// Puts the sheet's clipboard to the selected cell's row.
func (s *Sheet) PutRow() {
	row := s.SelectedCell.Row()

	rowSrc := s.clipboardRangeStart.Row()
	colStart := s.clipboardRangeStart.Column()
	colEnd := s.clipboardRangeEnd.Column()

	for colIdx := colStart; colIdx <= colEnd; colIdx++ {
		if srcCell, err := s.GetCell(NewAddress(rowSrc, colIdx)); err == nil {
			s.SetCell(NewAddress(row, colIdx), srcCell.Copy(NewAddress(rowSrc, colIdx), NewAddress(row, colIdx)))
		}

	}
}

// Puts the sheet's clipboard the the selected cell's column.
func (s *Sheet) PutColumn() {
	column := s.SelectedCell.Column()

	colSrc := s.clipboardRangeStart.Column()
	rowStart := s.clipboardRangeStart.Column()
	rowEnd := s.clipboardRangeEnd.Column()

	for rowIdx := rowStart; rowIdx <= rowEnd; rowIdx++ {
		if srcCell, err := s.GetCell(NewAddress(rowIdx, colSrc)); err == nil {
			s.SetCell(NewAddress(rowIdx, column), srcCell.Copy(NewAddress(rowIdx, colSrc), NewAddress(rowIdx, colSrc)))
		}

	}
}
