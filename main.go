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

type SheetMode int

const (
	NORMAL_MODE SheetMode = iota
	INSERT_MODE SheetMode = iota
)

func processTermboxEvents(s *Sheet) {
	smode := NORMAL_MODE
	valBuffer := bytes.Buffer{}
	for ev := termbox.PollEvent(); ev.Type != termbox.EventError; ev = termbox.PollEvent() {
		switch ev.Type {
		case termbox.EventKey:
			switch smode {
			case NORMAL_MODE:
				switch ev.Key {
				case termbox.KeyArrowUp:
					s.MoveUp()
				case termbox.KeyArrowDown:
					s.MoveDown()
				case termbox.KeyArrowLeft:
					s.MoveLeft()
				case termbox.KeyArrowRight:
					s.MoveRight()
				case 0:
					switch ev.Ch {
					case 'q':
						termbox.Close()
						return
					case '=', 'i':
						smode = INSERT_MODE
					case '<', '>', '|':
						smode = INSERT_MODE
						valBuffer.WriteRune(ev.Ch)
					case 'h':
						s.MoveLeft()
					case 'j':
						s.MoveDown()
					case 'k':
						s.MoveUp()
					case 'l':
						s.MoveRight()
					}
				}
			case INSERT_MODE:
				if ev.Key == termbox.KeyEnter {
					s.setCell(s.selectedCell, valBuffer.String())
					valBuffer.Reset()
					smode = NORMAL_MODE
				} else if ev.Key == termbox.KeyEsc {
					valBuffer.Reset()
					smode = NORMAL_MODE
				} else {
					valBuffer.WriteRune(ev.Ch)
				}
			}
		}
		displayValue(fmt.Sprintf("%s = %s", s.selectedCell, valBuffer.String()), 0, 0, 80, AlignLeft, false)
		termbox.Flush()
	}
}
