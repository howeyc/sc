package sheet

import (
	"fmt"

	"github.com/scrouthtv/gosc/internal/display"
	"github.com/scrouthtv/gosc/internal/sheet/align"

	"github.com/nsf/termbox-go"
)

const (
	DISPLAY_RAW_VALUE_ROW    = 0
	DISPLAY_COMMAND_HELP_ROW = 1
	DISPLAY_SHEET_START_ROW  = 2
)

// Display a window of the sheet using termbox.
//
// The top-left changes based on where in the sheet the selected cell is moved by movement commands.
func (s *Sheet) display() {
	if s.loading {
		return
	}
	defer termbox.Flush()
	displayWidth, displayHeight := termbox.Size()

	selectedRow, selectedCol := s.SelectedCell.RowCol()

	// Column Headers
	rowStr := fmt.Sprintf("%3d", s.startRow+displayHeight)
	x := 0
	for x <= len(rowStr) {
		termbox.SetCell(x, DISPLAY_SHEET_START_ROW, ' ', termbox.ColorWhite, termbox.ColorWhite)
		x++
	}
	startDispColumn := x
	displayColumns := 0
	columnAddr := NewAddress(0, s.startCol)
	for column := s.startCol; x+s.getColumnWidth(columnAddr.ColumnHeader()) < displayWidth; column++ {
		columnAddr = NewAddress(0, column)
		display.DisplayValue(columnAddr.ColumnHeader(), DISPLAY_SHEET_START_ROW, x, x+s.getColumnWidth(columnAddr.ColumnHeader()), align.AlignCenter, selectedCol != column)
		x += s.getColumnWidth(columnAddr.ColumnHeader())
		displayColumns = column - s.startCol + 1
	}

	displayRows := 0
	y := DISPLAY_SHEET_START_ROW + 1
	for row := s.startRow; y < displayHeight; y++ {
		rowStr := fmt.Sprintf("%3d", row)
		display.DisplayValue(rowStr, y, 0, len(rowStr)-1, align.AlignRight, selectedRow != row)
		displayRows = row - s.startRow + 1
		row++
	}

	termCol := startDispColumn
	for column := 0; column < displayColumns; column++ {
		valCol := column + s.startCol
		for row := 0; row < displayRows; row++ {
			valRow := row + s.startRow
			address := NewAddress(valRow, valCol)
			if cell, err := s.GetCell(address); err == nil {
				cell.display(s, address, row+DISPLAY_SHEET_START_ROW+1, termCol, termCol+s.getColumnWidth(address.ColumnHeader()), s.SelectedCell == address)
			} else {
				display.DisplayValue("", row+DISPLAY_SHEET_START_ROW+1, termCol, termCol+s.getColumnWidth(address.ColumnHeader()), align.AlignLeft, false)
			}
		}
		columnAddr := NewAddress(0, valCol)
		termCol += s.getColumnWidth(columnAddr.ColumnHeader())
	}
	s.displayRows, s.displayCols = displayRows, displayColumns
}
