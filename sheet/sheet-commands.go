// sheet-commands
package sheet

func (s *Sheet) YankRow() {
	row := s.SelectedCell.Row()

	s.clipboardRangeStart = getAddress(row, 0)
	s.clipboardRangeEnd = getAddress(row, s.maxColumnForRow[row])
}

func (s *Sheet) YankColumn() {
	column := s.SelectedCell.Column()

	s.clipboardRangeStart = getAddress(0, column)
	s.clipboardRangeEnd = getAddress(s.maxRowForColumn[column], column)
}

func (s *Sheet) PutRow() {
	row := s.SelectedCell.Row()

	rowSrc := s.clipboardRangeStart.Row()
	colStart := s.clipboardRangeStart.Column()
	colEnd := s.clipboardRangeEnd.Column()

	for colIdx := colStart; colIdx <= colEnd; colIdx++ {
		if srcCell, err := s.GetCell(getAddress(rowSrc, colIdx)); err == nil {
			s.SetCell(getAddress(row, colIdx), srcCell.Copy(getAddress(rowSrc, colIdx), getAddress(row, colIdx)))
		}

	}
}

func (s *Sheet) PutColumn() {
	column := s.SelectedCell.Column()

	colSrc := s.clipboardRangeStart.Column()
	rowStart := s.clipboardRangeStart.Column()
	rowEnd := s.clipboardRangeEnd.Column()

	for rowIdx := rowStart; rowIdx <= rowEnd; rowIdx++ {
		if srcCell, err := s.GetCell(getAddress(rowIdx, colSrc)); err == nil {
			s.SetCell(getAddress(rowIdx, column), srcCell.Copy(getAddress(rowIdx, colSrc), getAddress(rowIdx, colSrc)))
		}

	}
}
