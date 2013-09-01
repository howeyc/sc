package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

var columnArr []string

func init() {
	columnArr = make([]string, 26)
	columnArr[0] = "A"
	columnArr[1] = "B"
	columnArr[2] = "C"
	columnArr[3] = "D"
	columnArr[4] = "E"
	columnArr[5] = "F"
	columnArr[6] = "G"
	columnArr[7] = "H"
	columnArr[8] = "I"
	columnArr[9] = "J"
	columnArr[10] = "K"
	columnArr[11] = "L"
	columnArr[12] = "M"
	columnArr[13] = "N"
	columnArr[14] = "O"
	columnArr[15] = "P"
	columnArr[16] = "Q"
	columnArr[17] = "R"
	columnArr[18] = "S"
	columnArr[19] = "T"
}

type Cell struct {
	rawVal, dispVal string

	forwardRefs map[string]struct{} // Cells that are required for any formula
	backRefs    map[string]struct{} // Cells that reference this cell's value
}

func (c *Cell) display(row, colStart, colEnd int, selected bool) {
	displayValue(c.dispVal, row, colStart, colEnd, AlignRight, selected)
}

const (
	STARTING_COLUMN_WIDTH = 10
)

type Sheet struct {
	filename string

	selectedCell string
	columnWidths map[string]int
	data         map[string]*Cell
}

func newSheet(filename string) Sheet {
	s := Sheet{filename: filename, selectedCell: "A0", columnWidths: make(map[string]int), data: make(map[string]*Cell)}
	return s
}

func (s *Sheet) getColumnWidth(column string) int {
	if width, found := s.columnWidths[column]; found {
		return width
	} else {
		s.columnWidths[column] = STARTING_COLUMN_WIDTH
		return STARTING_COLUMN_WIDTH
	}
}

func (s *Sheet) increaseColumnWidth(column string) {
	if width, found := s.columnWidths[column]; found {
		s.columnWidths[column] = width + 1
	} else {
		s.columnWidths[column] = STARTING_COLUMN_WIDTH + 1
	}
}

func (s *Sheet) decreaseColumnWidth(column string) {
	if width, found := s.columnWidths[column]; found {
		if width > 1 {
			s.columnWidths[column] = width - 1
		}
	} else {
		s.columnWidths[column] = STARTING_COLUMN_WIDTH - 1
	}
}

func (s *Sheet) getCell(address string) (*Cell, error) {
	if cell, found := s.data[address]; found {
		return cell, nil
	}
	return nil, errors.New("Cell does not exist in spreadsheet.")
}

func (s *Sheet) setCell(address, val string) {
	// TODO: more work here to set refs and calc disp value
	s.data[address] = &Cell{rawVal: val, dispVal: val}
	s.display(0, 0)
}

type Align int

const (
	AlignRight Align = iota
	AlignLeft
	AlignCenter
)

func displayValue(val string, row, colStart, colEnd int, alignment Align, inverse bool) {
	fg, bg := termbox.ColorWhite, termbox.ColorBlack
	if inverse {
		fg, bg = bg, fg
	}
	valLen := utf8.RuneCountInString(val)
	rr := strings.NewReader(val)
	colWidth := colEnd - colStart + 1
	blankSize := colWidth - valLen
	if blankSize < 0 {
		blankSize = 0
	}
	startBlank, endBlank := 0, 0
	switch alignment {
	case AlignRight:
		startBlank = blankSize
	case AlignCenter:
		startBlank, endBlank = blankSize/2, blankSize/2
		if startBlank+endBlank < blankSize {
			endBlank++
		}
	case AlignLeft:
		endBlank = blankSize
	}
	i := 0
	for bs := 0; bs < startBlank; bs++ {
		termbox.SetCell(colStart+i, row, ' ', bg, bg)
		i++
	}
	runeSize := valLen
	if valLen > colWidth {
		runeSize = colWidth
	}
	for ri := 0; ri < runeSize; ri++ {
		nr, _, _ := rr.ReadRune()
		termbox.SetCell(colStart+i, row, nr, fg, bg)
		i++
	}
	for bs := 0; bs < endBlank; bs++ {
		termbox.SetCell(colStart+i, row, ' ', bg, bg)
		i++
	}
}

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
			}
		}
		termCol += s.getColumnWidth(columnArr[column])
	}
}
