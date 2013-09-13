// cell
package sheet

import (
	"fmt"
	"io"

	"scim/display"
	"scim/evaler"
	"scim/sheet/align"
)

type Cell struct {
	value      string
	stringType bool
	alignment  align.Align

	forwardRefs map[string]struct{} // Cells that are required for any formula
	backRefs    map[string]struct{} // Cells that reference this cell's value
}

func NewCell(value string, alignment align.Align, stringType bool) *Cell {
	return &Cell{value: value, alignment: alignment, stringType: stringType}
}

func (c *Cell) getDisplayValue(s *Sheet, address string) string {
	postfix := evaler.GetPostfix(c.value)
	for idx, token := range postfix {
		if evaler.IsCellAddr(token) {
			refCellVal := ""
			if cell, err := s.GetCell(token); err == nil {
				refCellVal = cell.getDisplayValue(s, token)
			}
			postfix[idx] = refCellVal
		}
	}
	if rat, err := evaler.EvaluatePostfix(postfix); err == nil {
		return rat.FloatString(s.getPrecision(address))
	} else {
		return ""
	}
}

func (c *Cell) display(s *Sheet, address string, row, colStart, colEnd int, selected bool) {
	dispVal := c.value
	if !c.stringType {
		dispVal = c.getDisplayValue(s, address)
	}
	display.DisplayValue(dispVal, row, colStart, colEnd, c.alignment, selected)
}

func (c *Cell) StatusBarVal() string {
	if c.stringType {
		modifier := ""
		switch c.alignment {
		case align.AlignLeft:
			modifier = "<"
		case align.AlignCenter:
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
	if c.stringType && c.alignment == align.AlignLeft {
		fmt.Fprintf(w, "leftstring %s = \"%s\"\n", address, c.value)
	} else if c.stringType && c.alignment == align.AlignCenter {
		fmt.Fprintf(w, "label %s = \"%s\"\n", address, c.value)
	} else if c.stringType {
		fmt.Fprintf(w, "rightstring %s = \"%s\"\n", address, c.value)
	} else {
		fmt.Fprintf(w, "let %s = %s\n", address, c.value)
	}
}
