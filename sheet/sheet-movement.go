// sheet-movement
package sheet

import (
	"fmt"
	"strconv"
)

func (s *Sheet) MoveUp() {
	sel := s.SelectedCell
	colStr := sel[:1]
	rowStr := sel[1:]
	row, _ := strconv.ParseInt(rowStr, 10, 64)
	row--
	if row < 0 {
		row = 0
	}
	if int(row) < s.startRow {
		s.startRow = int(row)
	}
	s.SelectedCell = fmt.Sprintf("%s%d", colStr, row)
	s.display()
}

func (s *Sheet) MoveDown() {
	sel := s.SelectedCell
	colStr := sel[:1]
	rowStr := sel[1:]
	row, _ := strconv.ParseInt(rowStr, 10, 64)
	row++
	s.SelectedCell = fmt.Sprintf("%s%d", colStr, row)
	if (int(row)-s.startRow)+1 > s.displayRows {
		s.startRow++
	}
	s.display()
}

func (s *Sheet) MoveRight() {
	sel := s.SelectedCell
	colStr := sel[:1]
	rowStr := sel[1:]
	row, _ := strconv.ParseInt(rowStr, 10, 64)
	colIdx := 0
	for columnArr[colIdx] != colStr {
		colIdx++
	}
	colIdx++
	if (colIdx-s.startCol)+1 > s.displayCols {
		s.startCol++
	}
	colStr = columnArr[colIdx]
	s.SelectedCell = fmt.Sprintf("%s%d", colStr, row)
	s.display()
}

func (s *Sheet) MoveLeft() {
	sel := s.SelectedCell
	colStr := sel[:1]
	rowStr := sel[1:]
	row, _ := strconv.ParseInt(rowStr, 10, 64)
	colIdx := 0
	for columnArr[colIdx] != colStr {
		colIdx++
	}
	colIdx--
	if colIdx < 0 {
		colIdx = 0
	}
	if colIdx < s.startCol {
		s.startCol = colIdx
	}
	colStr = columnArr[colIdx]
	s.SelectedCell = fmt.Sprintf("%s%d", colStr, row)
	s.display()
}
