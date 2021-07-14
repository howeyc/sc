// Implements functions required commands on a sheet.
package sheet

import (
	"errors"
	"fmt"

	"github.com/scrouthtv/gosc/internal/evaler"
)

const (
	initial_COLUMN_WIDTH     = 10
	initial_COLUMN_PRECISION = 2
	initial_COLUMN_TYPE      = 0
)

type ColumnFormat struct {
	width     int
	precision int
	ctype     int
}

type Sheet struct {
	Filename string

	SelectedCell  Address
	columnFormats map[string]ColumnFormat
	data          map[Address]*Cell

	maxRowForColumn map[int]int
	maxColumnForRow map[int]int

	clipboardRangeStart, clipboardRangeEnd Address

	// display window
	startRow, startCol       int
	displayRows, displayCols int

	loading bool
}

// Creates a new sheet and loads a sheet from filename if it exists.
func NewSheet(filename string) Sheet {
	s := Sheet{Filename: filename, SelectedCell: "A0",
		columnFormats: make(map[string]ColumnFormat), data: make(map[Address]*Cell),
		maxRowForColumn: make(map[int]int), maxColumnForRow: make(map[int]int)}

	// Run display to get the row and column height of display set inside of sheet.
	s.display()

	// Load values into sheet and move to selected cell
	s.Load()

	// Display loaded sheet
	s.display()

	return s
}

// Given an address, returns the precision to use to display a number.
func (s *Sheet) getPrecision(address Address) int {
	if cFormat, found := s.columnFormats[address.ColumnHeader()]; found {
		return cFormat.precision
	} else {
		return initial_COLUMN_PRECISION
	}
}

// For a given column header string, decreases column precision.
func (s *Sheet) DecreaseColumnPrecision(column string) {
	if cFormat, found := s.columnFormats[column]; found {
		cFormat.precision--
		if cFormat.precision < 0 {
			cFormat.precision = 0
		}
		s.columnFormats[column] = cFormat
	} else {
		s.columnFormats[column] = ColumnFormat{width: initial_COLUMN_WIDTH, precision: initial_COLUMN_PRECISION - 1, ctype: initial_COLUMN_TYPE}
	}
	s.display()
}

// For a given column header string, increases column precision.
func (s *Sheet) IncreaseColumnPrecision(column string) {
	if cFormat, found := s.columnFormats[column]; found {
		cFormat.precision++
		s.columnFormats[column] = cFormat
	} else {
		s.columnFormats[column] = ColumnFormat{width: initial_COLUMN_WIDTH, precision: initial_COLUMN_PRECISION + 1, ctype: initial_COLUMN_TYPE}
	}
	s.display()
}

// Returns the the display format string specifying the Width, Precition, and Type
// of the column.
func (s *Sheet) DisplayFormat(address Address) string {
	if cFormat, found := s.columnFormats[address.ColumnHeader()]; found {
		return fmt.Sprintf("%d %d %d", cFormat.width, cFormat.precision, cFormat.ctype)
	} else {
		return fmt.Sprintf("%d %d %d", initial_COLUMN_WIDTH, initial_COLUMN_PRECISION, initial_COLUMN_TYPE)
	}
}

// For a given column header string, returns the display width of the column in characters.
func (s *Sheet) getColumnWidth(column string) int {
	if cFormat, found := s.columnFormats[column]; found {
		return cFormat.width
	} else {
		s.columnFormats[column] = ColumnFormat{width: initial_COLUMN_WIDTH, precision: initial_COLUMN_PRECISION, ctype: initial_COLUMN_TYPE}
		return initial_COLUMN_WIDTH
	}
}

// For a given column header string, decreases column width.
func (s *Sheet) DecreaseColumnWidth(column string) {
	if cFormat, found := s.columnFormats[column]; found {
		cFormat.width--
		if cFormat.width < 1 {
			cFormat.width = 1
		}
		s.columnFormats[column] = cFormat
	} else {
		s.columnFormats[column] = ColumnFormat{width: initial_COLUMN_WIDTH - 1, precision: initial_COLUMN_PRECISION, ctype: initial_COLUMN_TYPE}
	}
	s.display()
}

// For a given column header string, increases column width.
func (s *Sheet) IncreaseColumnWidth(column string) {
	if cFormat, found := s.columnFormats[column]; found {
		cFormat.width++
		s.columnFormats[column] = cFormat
	} else {
		s.columnFormats[column] = ColumnFormat{width: initial_COLUMN_WIDTH + 1, precision: initial_COLUMN_PRECISION, ctype: initial_COLUMN_TYPE}
	}
	s.display()
}

// Removes the cell at the given address from a sheet.
func (s *Sheet) ClearCell(address Address) {
	delete(s.data, address)

	s.findMaximums(address)
}

// Returns the cell at the given address.
func (s *Sheet) GetCell(address Address) (*Cell, error) {
	if cell, found := s.data[address]; found {
		return cell, nil
	} else if address == s.SelectedCell {
		return &Cell{}, nil
	}
	return nil, errors.New("Cell does not exist in spreadsheet.")
}

// Sets the address to the passed in cell. Previous cell data that exists is thrown away.
func (s *Sheet) SetCell(address Address, cell *Cell) {
	if !cell.stringType {
		postfix := evaler.GetPostfix(cell.value)
		for _, token := range postfix {
			if evaler.IsCellAddr(token) {
				tokenAddr := Address(token)
				cell.forwardRefs[tokenAddr] = struct{}{}
			}
		}
	}
	s.data[address] = cell

	s.setMaximums(address)

	// Refresh the sheet
	s.display()
}

// Sets the max row and column with values in them to make clipboard actions possible.
func (s *Sheet) setMaximums(address Address) {
	row, col := address.RowCol()

	currRowMax := s.maxRowForColumn[col]
	if row > currRowMax {
		s.maxRowForColumn[col] = row
	}

	currColMax := s.maxColumnForRow[row]
	if col > currColMax {
		s.maxColumnForRow[row] = col
	}
}

// Finds the max row and column with values in them to make clipboard actions possible.
func (s *Sheet) findMaximums(address Address) {
	row, column := address.RowCol()

	currColMax := s.maxColumnForRow[row]
	if column == currColMax {
		// Find lower column used on row
		delete(s.maxColumnForRow, row)
		for colIdx := column; colIdx >= 0; colIdx-- {
			if _, found := s.data[NewAddress(row, colIdx)]; found {
				s.maxColumnForRow[row] = colIdx
			}
		}
	}

	currRowMax := s.maxRowForColumn[column]
	if row == currRowMax {
		// Find lower row used for column
		delete(s.maxRowForColumn, column)
		for rowIdx := row; rowIdx >= 0; rowIdx-- {
			if _, found := s.data[NewAddress(rowIdx, column)]; found {
				s.maxRowForColumn[column] = rowIdx
			}
		}
	}
}
