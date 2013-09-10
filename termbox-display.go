// termbox-display
package main

import (
	"strings"
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

type Align int

const (
	AlignRight Align = iota
	AlignLeft
	AlignCenter
)

func displayValue(val string, row, colStart, colEnd int, alignment Align, inverse bool) {
	fg, bg := termbox.ColorWhite, termbox.ColorBlack
	if inverse {
		fg, bg = bg, fg
	}
	valLen := utf8.RuneCountInString(val)
	rr := strings.NewReader(val)
	colWidth := colEnd - colStart + 1
	blankSize := colWidth - valLen
	if blankSize < 0 {
		blankSize = 0
	}
	startBlank, endBlank := 0, 0
	switch alignment {
	case AlignRight:
		startBlank = blankSize - 1
	case AlignCenter:
		startBlank, endBlank = blankSize/2, blankSize/2
		if startBlank+endBlank < blankSize {
			endBlank++
		}
	case AlignLeft:
		endBlank = blankSize
	}
	i := 0
	for bsl := 0; bsl < startBlank; bsl++ {
		termbox.SetCell(colStart+i, row, ' ', bg, bg)
		i++
	}
	runeSize := valLen
	if valLen > colWidth {
		runeSize = colWidth
	}
	for ri := 0; ri < runeSize; ri++ {
		nr, _, _ := rr.ReadRune()
		termbox.SetCell(colStart+i, row, nr, fg, bg)
		i++
	}
	for bsr := 0; bsr < endBlank; bsr++ {
		termbox.SetCell(colStart+i, row, ' ', bg, bg)
		i++
	}
}
