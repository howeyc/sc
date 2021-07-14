package sheet

import (
	"fmt"
	"strings"

	"github.com/scrouthtv/gosc/internal/display"
	"github.com/scrouthtv/gosc/internal/evaler"
	"github.com/scrouthtv/gosc/internal/sheet/align"
)

type Cell struct {
	value      string
	stringType bool
	alignment  align.Align

	forwardRefs map[Address]struct{} // Cells that are required for any formula
}

// Creates a new cell.
func NewCell(value string, alignment align.Align, stringType bool) *Cell {
	return &Cell{
		value:       value,
		alignment:   alignment,
		stringType:  stringType,
		forwardRefs: make(map[Address]struct{}),
	}
}

// Creates a copy of the cell, altering any formula that is contained based on where the
// cell is being moved. It uses the relative change between the old address and new address
// to figure out the differences to apply to every cell referenced in the formula.
func (c *Cell) Copy(oldAddress, newAddress Address) *Cell {
	val := c.value
	if !c.stringType {
		rowDiff := newAddress.Row() - oldAddress.Row()
		colDiff := newAddress.Column() - oldAddress.Column()
		for forRef, _ := range c.forwardRefs {
			refRow, refCol := forRef.RowCol()
			newRef := NewAddress(refRow+rowDiff, refCol+colDiff)
			val = strings.Replace(val, string(forRef), string(newRef), -1)
		}
	}
	return NewCell(val, c.alignment, c.stringType)
}

// Returns the value to display as a string. This is the value that shows up on
// cell in the sheet display. This is the result of any calculation that may have
// been required to be performed.
//
// Cell is assumed not to be string type when this function is called.
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

// Displays cell value using termbox.
func (c *Cell) display(s *Sheet, address Address, row, colStart, colEnd int, selected bool) {
	dispVal := c.value
	if !c.stringType {
		dispVal = c.getDisplayValue(s, address)
	}
	display.DisplayValue(dispVal, row, colStart, colEnd, c.alignment, selected)
}

// Gets the raw value in a format the also specifies any alignment defined in cell.
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
