// sheet-display
package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

const (
	DISPLAY_RAW_VALUE_ROW    = 0
	DISPLAY_COMMAND_HELP_ROW = 1
	DISPLAY_SHEET_START_ROW  = 2
)

func (s *Sheet) display(startRow, startColumn int) {
	defer termbox.Flush()
	displayWidth, displayHeight := termbox.Size()

	// Column Headers
	rowStr := fmt.Sprintf("% 3d", startRow)
	x := 0
	for x <= len(rowStr) {
		termbox.SetCell(x, DISPLAY_SHEET_START_ROW, ' ', termbox.ColorWhite, termbox.ColorWhite)
		x++
	}
	startDispColumn := x
	displayColumns := 0
	for column := 0; x+s.getColumnWidth(columnArr[column]) < displayWidth; column++ {
		displayValue(columnArr[column], DISPLAY_SHEET_START_ROW, x, x+s.getColumnWidth(columnArr[column]), AlignCenter, true)
		x += s.getColumnWidth(columnArr[column])
		displayColumns = column
	}

	displayRows := 0
	for row := DISPLAY_SHEET_START_ROW + 1; row < displayHeight; row++ {
		rowStr := fmt.Sprintf("% 3d", row-(DISPLAY_SHEET_START_ROW+1))
		displayValue(rowStr, row, 0, len(rowStr)-1, AlignRight, true)
		displayRows = row
	}
	displayRows = displayRows - (DISPLAY_SHEET_START_ROW + 1)

	termCol := startDispColumn
	for column := 0; column < displayColumns; column++ {
		for row := 0; row < displayRows; row++ {
			address := fmt.Sprintf("%s%d", columnArr[column], row)
			if cell, err := s.getCell(address); err == nil {
				cell.display(row+DISPLAY_SHEET_START_ROW+1, termCol, termCol+s.getColumnWidth(columnArr[column]), s.selectedCell == address)
			} else {
				displayValue("", row+DISPLAY_SHEET_START_ROW+1, termCol, termCol+s.getColumnWidth(columnArr[column]), AlignLeft, false)
			}
		}
		termCol += s.getColumnWidth(columnArr[column])
	}
}
