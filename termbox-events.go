// termbox-events
package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

type SheetMode int

const (
	NORMAL_MODE SheetMode = iota
	INSERT_MODE SheetMode = iota
	EXIT_MODE   SheetMode = iota
)

func processTermboxEvents(s *Sheet) {
	prompt := ""
	stringEntry := false
	smode := NORMAL_MODE
	valBuffer := bytes.Buffer{}

	// Display
	go func() {
		for _ = range time.Tick(200 * time.Millisecond) {
			switch smode {
			case NORMAL_MODE:
				selSel, _ := s.getCell(s.selectedCell)
				displayValue(fmt.Sprintf("%s (10 2 0) [%s]", s.selectedCell, selSel.rawVal), 0, 0, 80, AlignLeft, false)
			case INSERT_MODE:
				displayValue(fmt.Sprintf("i> %s %s = %s", prompt, s.selectedCell, valBuffer.String()), 0, 0, 80, AlignLeft, false)
			case EXIT_MODE:
				displayValue(fmt.Sprintf("File \"%s\" is modified, save before exiting?", "<file>"), 0, 0, 80, AlignLeft, false)
			}
			termbox.Flush()
		}
	}()

	// Events
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
						smode = EXIT_MODE
					case '=', 'i':
						smode = INSERT_MODE
						prompt = "let"
					case '<':
						prompt = "leftstring"
						smode = INSERT_MODE
						valBuffer.WriteRune(ev.Ch)
						valBuffer.WriteRune('"')
						stringEntry = true
					case '>':
						prompt = "rightstring"
						smode = INSERT_MODE
						valBuffer.WriteRune(ev.Ch)
						valBuffer.WriteRune('"')
						stringEntry = true
					case '\\':
						prompt = "label"
						smode = INSERT_MODE
						valBuffer.WriteRune('|')
						valBuffer.WriteRune('"')
						stringEntry = true
					case 'h':
						s.MoveLeft()
					case 'j':
						s.MoveDown()
					case 'k':
						s.MoveUp()
					case 'l':
						s.MoveRight()
					case 'x':
						s.clearCell(s.selectedCell)
					}
				}
			case INSERT_MODE:
				if ev.Key == termbox.KeyEnter {
					if stringEntry {
						valBuffer.WriteRune('"')
					}
					s.setCell(s.selectedCell, valBuffer.String())
					valBuffer.Reset()
					smode = NORMAL_MODE
					stringEntry = false
				} else if ev.Key == termbox.KeyEsc {
					valBuffer.Reset()
					smode = NORMAL_MODE
					stringEntry = false
				} else if ev.Key == termbox.KeyBackspace {
					s := valBuffer.String()
					valBuffer = bytes.Buffer{}
					if len(s) > 0 {
						s = s[0 : len(s)-1]
					}
					valBuffer.WriteString(s)
				} else {
					valBuffer.WriteRune(ev.Ch)
				}
			case EXIT_MODE:
				if ev.Key == 0 && ev.Ch == 'y' {
					// TODO: SAVE
				}
				termbox.Close()
				return
			}
		}
	}
}
