// cell
package main

import (
	"fmt"
	"io"
)

type Cell struct {
	rawVal, dispVal string
	alignment       Align

	forwardRefs map[string]struct{} // Cells that are required for any formula
	backRefs    map[string]struct{} // Cells that reference this cell's value
}

func NewCell(val string) *Cell {
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
	return &Cell{rawVal: val, dispVal: dispVal, alignment: alignment}
}

func (c *Cell) display(row, colStart, colEnd int, selected bool) {
	displayValue(c.dispVal, row, colStart, colEnd, c.alignment, selected)
}

func (c *Cell) write(w io.Writer, address string) {
	switch c.rawVal[0] {
	case '<':
		fmt.Fprintf(w, "leftstring %s = \"%s\"\n", address, c.dispVal)
	case '>':
		fmt.Fprintf(w, "rightstring %s = \"%s\"\n", address, c.dispVal)
	case '|':
		fmt.Fprintf(w, "label %s = \"%s\"\n", address, c.dispVal)
	default:
		fmt.Fprintf(w, "let %s = %s\n", address, c.dispVal)
	}
}
