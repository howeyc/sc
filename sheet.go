package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

func splitCellAddress(address string) (column, row int) {
	strCol := address[0]
	strRow := address[1:]

	column = int(strCol - 'A')
	rowI64, _ := strconv.ParseInt(strRow, 10, 64)
	row = int(rowI64)
	return column, row
}

type Cell struct {
	rawVal, dispVal string

	forwardRefs map[string][]string // Cells that are required for any formula
	backRefs    map[string][]string // Cells that reference this cell's value
}

func (c *Cell) display(row, colStart, colEnd int) {

}

type Sheet struct {
	filename string

	// Change to map[string]Cell ??
	columns [][]Cell
}

func (s *Sheet) getCell(address string) (*Cell, error) {
	column, row := splitCellAddress(address)
	if column < len(s.columns) && row < len(s.columns[column]) {
		return &s.columns[column][row], nil
	}
	return nil, errors.New("Cell does not exist in spreadsheet.")
}

func (s *Sheet) setCell(address, val string) {
	column, row := splitCellAddress(address)
	if column >= len(s.columns) {
		s.columns = append(s.columns, make([][]Cell, column-len(s.columns)+1)...)
	}
	if row >= len(s.columns[column]) {
		s.columns[column] = append(s.columns[column], make([]Cell, row-len(s.columns[column])+1)...)
	}

	// TODO: more work here to set refs and calc disp value
	s.columns[column][row] = Cell{rawVal: val}
}

func (s *Sheet) display() {
	defer termbox.Flush()

	// Column Headers
	x := 0
	for x < 4 {
		termbox.SetCell(x, 3, ' ', termbox.ColorWhite, termbox.ColorWhite)
		x++
	}
	for column := 0; column < 11; column++ {
		for fld := 0; fld < 10; fld++ {
			if fld == 4 {
				columnHeader := []byte{byte(65 + column)}
				runeHeader, _ := utf8.DecodeRune(columnHeader)
				termbox.SetCell(x, 3, runeHeader, termbox.ColorBlack, termbox.ColorWhite)
			} else {
				termbox.SetCell(x, 3, ' ', termbox.ColorWhite, termbox.ColorWhite)
			}
			x++
		}
	}

	for row := 4; row < 65; row++ {
		rowStr := fmt.Sprintf("% 3d", row-4)
		sr := strings.NewReader(rowStr)
		nr, _, _ := sr.ReadRune()
		termbox.SetCell(0, row, nr, termbox.ColorBlack, termbox.ColorWhite)
		nr, _, _ = sr.ReadRune()
		termbox.SetCell(1, row, nr, termbox.ColorBlack, termbox.ColorWhite)
		nr, _, _ = sr.ReadRune()
		termbox.SetCell(2, row, nr, termbox.ColorBlack, termbox.ColorWhite)
	}
}
