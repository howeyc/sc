package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		fmt.Println("Could not start termbox.")
		os.Exit(1)
	}

	sheet := newSheet("")

	// Set some default values
	sheet.setCell("A0", "start")
	sheet.setCell("A1", "adsf")
	sheet.setCell("A2", "ljl;")
	sheet.setCell("A3", "owerjf")
	sheet.setCell("A4", "woerjlfj")
	sheet.setCell("D4", "Roar")

	sheet.display(0, 0)

	processTermboxEvents(&sheet)
}

func processTermboxEvents(s *Sheet) {
	valBuffer := bytes.Buffer{}
	for ev := termbox.PollEvent(); ev.Type != termbox.EventError; ev = termbox.PollEvent() {
		if ev.Type == termbox.EventKey {
			if ev.Ch == 'q' {
				termbox.Close()
				return
			}
			if ev.Key == termbox.KeyArrowUp {
				s.MoveUp()
				continue
			}
			if ev.Key == termbox.KeyArrowDown {
				s.MoveDown()
				continue
			}
			if ev.Key == termbox.KeyEnter {
				s.setCell(s.selectedCell, valBuffer.String())
				valBuffer.Reset()
				continue
			}
			valBuffer.WriteRune(ev.Ch)
		}
		displayValue(fmt.Sprintf("%s = %s", s.selectedCell, valBuffer.String()), 0, 0, 80, AlignLeft, false)
		termbox.Flush()
	}
}
