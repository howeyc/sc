// cell
package main

import (
	"fmt"
	"io"
)

type Cell struct {
	value      string
	stringType bool
	alignment  Align

	forwardRefs map[string]struct{} // Cells that are required for any formula
	backRefs    map[string]struct{} // Cells that reference this cell's value
}

func (c *Cell) display(row, colStart, colEnd int, selected bool) {
	displayValue(c.value, row, colStart, colEnd, c.alignment, selected)
}

func (c *Cell) statusBarVal() string {
	if c.stringType {
		modifier := ""
		switch c.alignment {
		case AlignLeft:
			modifier = "<"
		case AlignCenter:
			modifier = "|"
		default:
			modifier = ">"
		}
		return fmt.Sprintf("%s\"%s\"", modifier, c.value)
	} else {
		return c.value
	}
}

func (c *Cell) write(w io.Writer, address string) {
	if c.stringType && c.alignment == AlignLeft {
		fmt.Fprintf(w, "leftstring %s = \"%s\"\n", address, c.value)
	} else if c.stringType && c.alignment == AlignCenter {
		fmt.Fprintf(w, "label %s = \"%s\"\n", address, c.value)
	} else if c.stringType {
		fmt.Fprintf(w, "rightstring %s = \"%s\"\n", address, c.value)
	} else {
		fmt.Fprintf(w, "let %s = %s\n", address, c.value)
	}
}
