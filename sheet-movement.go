// sheet-movement
package main

import (
	"fmt"
	"strconv"
)

func (s *Sheet) MoveUp() {
	sel := s.selectedCell
	colStr := sel[:1]
	rowStr := sel[1:]
	row, _ := strconv.ParseInt(rowStr, 10, 64)
	row--
	if row < 0 {
		row = 0
	}
	s.selectedCell = fmt.Sprintf("%s%d", colStr, row)
	s.display(0, 0)
}

func (s *Sheet) MoveDown() {
	sel := s.selectedCell
	colStr := sel[:1]
	rowStr := sel[1:]
	row, _ := strconv.ParseInt(rowStr, 10, 64)
	row++
	/* TODO: move sheet down if necessary
	if row < 0 {
		row = 0
	}
	*/
	s.selectedCell = fmt.Sprintf("%s%d", colStr, row)
	s.display(0, 0)
}

func (s *Sheet) MoveRight() {
	sel := s.selectedCell
	colStr := sel[:1]
	rowStr := sel[1:]
	row, _ := strconv.ParseInt(rowStr, 10, 64)
	colIdx := 0
	for columnArr[colIdx] != colStr {
		colIdx++
	}
	colIdx++
	colStr = columnArr[colIdx]
	s.selectedCell = fmt.Sprintf("%s%d", colStr, row)
	s.display(0, 0)
}

func (s *Sheet) MoveLeft() {
	sel := s.selectedCell
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
	colStr = columnArr[colIdx]
	s.selectedCell = fmt.Sprintf("%s%d", colStr, row)
	s.display(0, 0)
}
