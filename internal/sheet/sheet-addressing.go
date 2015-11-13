package sheet

import (
	"regexp"
	"strconv"
	"strings"
)

type Address string

var re_addr = regexp.MustCompile(`([A-Z]+)(\d+)`)

// Splits an address into row and column strings.
func (addr Address) Split() (row, column string) {
	parts := re_addr.FindStringSubmatch(string(addr))

	if len(parts) == 3 {
		column = parts[1]
		row = parts[2]
	}
	return
}

// Returns the integer value of the given address.
//
// Columns are zero-based, using the letters A-Z, starting with A = 0.
func (addr Address) Column() int {
	_, column := addr.Split()

	col := 0
	for len(column) > 0 {
		colNum, _ := strconv.ParseInt(column[:1], 36, 64)
		colNum -= 9
		col *= 26
		col += int(colNum)
		column = column[1:]
	}
	return int(col) - 1
}

// Returns the column portion of the address as a string.
func (addr Address) ColumnHeader() string {
	_, column := addr.Split()
	return column
}

// Returns the row portion of the address as an integer.
func (addr Address) Row() int {
	rowStr, _ := addr.Split()

	row, _ := strconv.ParseInt(rowStr, 10, 64)
	return int(row)
}

// Returns the row and column of the address as integers.
func (addr Address) RowCol() (row, column int) {
	return addr.Row(), addr.Column()
}

// Creates and address given an integer row and column.
func NewAddress(row, column int) Address {
	column += 10
	return Address(strings.ToUpper(strconv.FormatInt(int64(column), 36)) + strconv.FormatInt(int64(row), 10))
}
