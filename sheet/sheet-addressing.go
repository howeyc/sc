// sheet-addressing
package sheet

import (
	"regexp"
	"strconv"
	"strings"
)

type Address string

var re_addr = regexp.MustCompile(`([A-Z]+)(\d+)`)

func (addr Address) Split() (row, column string) {
	parts := re_addr.FindStringSubmatch(string(addr))

	if len(parts) == 3 {
		column = parts[1]
		row = parts[2]
	}
	return
}

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

func (addr Address) ColumnHeader() string {
	_, column := addr.Split()
	return column
}

func (addr Address) Row() int {
	rowStr, _ := addr.Split()

	row, _ := strconv.ParseInt(rowStr, 10, 64)
	return int(row)
}

func (addr Address) RowCol() (row, column int) {
	return addr.Row(), addr.Column()
}

func getAddress(row, column int) Address {
	column += 10
	return Address(strings.ToUpper(strconv.FormatInt(int64(column), 36)) + strconv.FormatInt(int64(row), 10))
}
