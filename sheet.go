package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
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

type ColumnFormat struct {
	width     int
	precision int
	ctype     int
}

type Sheet struct {
	filename string

	selectedCell  string
	columnFormats map[string]ColumnFormat
	data          map[string]*Cell

	// display window
	startRow, startCol       int
	displayRows, displayCols int
}

func newSheet(filename string) Sheet {
	s := Sheet{filename: filename, selectedCell: "A0", columnFormats: make(map[string]ColumnFormat), data: make(map[string]*Cell)}

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
			if len(val) > 1 && val[0] == '"' {
				val = val[1 : len(val)-1]
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
			case "format":
				width, _ := strconv.ParseInt(words[2], 10, 64)
				precision, _ := strconv.ParseInt(words[3], 10, 64)
				ctype, _ := strconv.ParseInt(words[4], 10, 64)
				s.columnFormats[adrs] = ColumnFormat{width: int(width), precision: int(precision), ctype: int(ctype)}
			}
		}
	}

	return s
}

func (s *Sheet) writeFormats(w io.Writer) {
	for k, cFormat := range s.columnFormats {
		fmt.Fprintf(w, "format %s %d %d %d\n", k, cFormat.width, cFormat.precision, cFormat.ctype)
	}
}

func (s *Sheet) displayFormat(address string) string {
	column := address[1:]
	if cFormat, found := s.columnFormats[column]; found {
		return fmt.Sprintf("%d %d %d", cFormat.width, cFormat.precision, cFormat.ctype)
	} else {
		return fmt.Sprintf("%d %d %d", STARTING_COLUMN_WIDTH, 2, 0)
	}
}

func (s *Sheet) getColumnWidth(column string) int {
	if cFormat, found := s.columnFormats[column]; found {
		return cFormat.width
	} else {
		s.columnFormats[column] = ColumnFormat{width: STARTING_COLUMN_WIDTH, precision: 2, ctype: 0}
		return STARTING_COLUMN_WIDTH
	}
}

func (s *Sheet) increaseColumnWidth(column string) {
	if cFormat, found := s.columnFormats[column]; found {
		cFormat.width += 1
		s.columnFormats[column] = cFormat
	} else {
		s.columnFormats[column] = ColumnFormat{width: STARTING_COLUMN_WIDTH + 1, precision: 2, ctype: 0}
	}
}

func (s *Sheet) decreaseColumnWidth(column string) {
	if cFormat, found := s.columnFormats[column]; found {
		if cFormat.width > 1 {
			cFormat.width--
			s.columnFormats[column] = cFormat
		}
	} else {
		s.columnFormats[column] = ColumnFormat{width: STARTING_COLUMN_WIDTH - 1, precision: 2, ctype: 0}
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
