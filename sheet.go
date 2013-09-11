package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
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

	// Load file
	if file, err := os.Open(filename); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") || len(line) < 1 {
				continue
			}
			words := strings.Split(line, " ")
			cmd := ""
			adrs := ""
			val := ""
			if len(words) >= 2 {
				cmd = words[0]
				adrs = words[1]
			}
			if len(words) >= 4 {
				val = strings.Join(words[3:], " ")
			}
			switch cmd {
			case "leftstring":
				s.setCell(adrs, &Cell{stringType: true, alignment: AlignLeft, value: val})
			case "rightstrng":
				s.setCell(adrs, &Cell{stringType: true, alignment: AlignRight, value: val})
			case "label":
				s.setCell(adrs, &Cell{stringType: true, alignment: AlignCenter, value: val})
			case "let":
				s.setCell(adrs, &Cell{stringType: false, alignment: AlignRight, value: val})
			case "goto":
				s.selectedCell = adrs
			}
		}
	}

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

func (s *Sheet) setCell(address string, cell *Cell) {
	// TODO: more work here to set refs and calc disp value
	s.data[address] = cell
	s.display()
}
