package main

import (
	"errors"
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

const (
	STARTING_COLUMN_WIDTH = 10
)

type Sheet struct {
	filename string

	selectedCell string
	columnWidths map[string]int
	data         map[string]*Cell

	// display window
	startRow, startCol       int
	displayRows, displayCols int
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

func (s *Sheet) clearCell(address string) {
	delete(s.data, address)
}

func (s *Sheet) getCell(address string) (*Cell, error) {
	if cell, found := s.data[address]; found {
		return cell, nil
	} else if address == s.selectedCell {
		return &Cell{}, nil
	}
	return nil, errors.New("Cell does not exist in spreadsheet.")
}

func (s *Sheet) setCell(address, val string) {
	// TODO: more work here to set refs and calc disp value
	alignment := AlignRight
	dispVal := val
	if val[0] == '<' {
		alignment = AlignLeft
		dispVal = val[2 : len(val)-1]
	} else if val[0] == '>' {
		alignment = AlignRight
		dispVal = val[2 : len(val)-1]
	} else if val[0] == '|' {
		alignment = AlignCenter
		dispVal = val[2 : len(val)-1]
	}
	s.data[address] = &Cell{rawVal: val, dispVal: dispVal, alignment: alignment}
	s.display()
}
