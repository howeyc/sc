// cell
package sheet

import (
	"fmt"
	"io"
	"strings"

	"github.com/howeyc/sc/display"
	"github.com/howeyc/sc/evaler"
	"github.com/howeyc/sc/sheet/align"
)

type Cell struct {
	value      string
	stringType bool
	alignment  align.Align

	forwardRefs map[Address]struct{} // Cells that are required for any formula
	backRefs    map[Address]struct{} // Cells that reference this cell's value
}

func NewCell(value string, alignment align.Align, stringType bool) *Cell {
	return &Cell{value: value, alignment: alignment, stringType: stringType,
		forwardRefs: make(map[Address]struct{}), backRefs: make(map[Address]struct{})}
}

func (c *Cell) Copy(oldAddress, newAddress Address) *Cell {
	val := c.value
	if !c.stringType {
		rowDiff := newAddress.Row() - oldAddress.Row()
		colDiff := newAddress.Column() - oldAddress.Column()
		for forRef, _ := range c.forwardRefs {
			refRow, refCol := forRef.RowCol()
			newRef := getAddress(refRow+rowDiff, refCol+colDiff)
			val = strings.Replace(val, string(forRef), string(newRef), -1)
		}
	}
	return NewCell(val, c.alignment, c.stringType)
}

func (c *Cell) getDisplayValue(s *Sheet, address Address) string {
	postfix := evaler.GetPostfix(c.value)
	for idx, token := range postfix {
		if evaler.IsCellAddr(token) {
			refCellVal := ""
			if cell, err := s.GetCell(Address(token)); err == nil {
				refCellVal = cell.getDisplayValue(s, Address(token))
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

func (c *Cell) display(s *Sheet, address Address, row, colStart, colEnd int, selected bool) {
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
