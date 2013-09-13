// sheet-display
package sheet

import (
	"fmt"

	"scim/display"
	"scim/sheet/align"

	"github.com/nsf/termbox-go"
)

const (
	DISPLAY_RAW_VALUE_ROW    = 0
	DISPLAY_COMMAND_HELP_ROW = 1
	DISPLAY_SHEET_START_ROW  = 2
)

func (s *Sheet) display() {
	defer termbox.Flush()
	displayWidth, displayHeight := termbox.Size()

	// Column Headers
	rowStr := fmt.Sprintf("% 3d", s.startRow)
	x := 0
	for x <= len(rowStr) {
		termbox.SetCell(x, DISPLAY_SHEET_START_ROW, ' ', termbox.ColorWhite, termbox.ColorWhite)
		x++
	}
	startDispColumn := x
	displayColumns := 0
	for column := s.startCol; x+s.getColumnWidth(columnArr[column]) < displayWidth; column++ {
		display.DisplayValue(columnArr[column], DISPLAY_SHEET_START_ROW, x, x+s.getColumnWidth(columnArr[column]), align.AlignCenter, true)
		x += s.getColumnWidth(columnArr[column])
		displayColumns = column - s.startCol + 1
	}

	displayRows := 0
	y := DISPLAY_SHEET_START_ROW + 1
	for row := s.startRow; y < displayHeight; y++ {
		rowStr := fmt.Sprintf("% 3d", row)
		display.DisplayValue(rowStr, y, 0, len(rowStr)-1, align.AlignRight, true)
		displayRows = row - s.startRow + 1
		row++
	}

	termCol := startDispColumn
	for column := 0; column < displayColumns; column++ {
		valCol := column + s.startCol
		for row := 0; row < displayRows; row++ {
			valRow := row + s.startRow
			address := fmt.Sprintf("%s%d", columnArr[valCol], valRow)
			if cell, err := s.GetCell(address); err == nil {
				cell.display(s, address, row+DISPLAY_SHEET_START_ROW+1, termCol, termCol+s.getColumnWidth(columnArr[valCol]), s.SelectedCell == address)
			} else {
				display.DisplayValue("", row+DISPLAY_SHEET_START_ROW+1, termCol, termCol+s.getColumnWidth(columnArr[valCol]), align.AlignLeft, false)
			}
		}
		termCol += s.getColumnWidth(columnArr[valCol])
	}
	s.displayRows, s.displayCols = displayRows, displayColumns
}
