// cell
package main

type Cell struct {
	rawVal, dispVal string
	alignment       Align

	forwardRefs map[string]struct{} // Cells that are required for any formula
	backRefs    map[string]struct{} // Cells that reference this cell's value
}

func (c *Cell) display(row, colStart, colEnd int, selected bool) {
	displayValue(c.dispVal, row, colStart, colEnd, c.alignment, selected)
}
