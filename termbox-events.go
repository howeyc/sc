// termbox-events
package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/howeyc/sc/display"
	"github.com/howeyc/sc/sheet"
	"github.com/howeyc/sc/sheet/align"
	"github.com/nsf/termbox-go"
)

type SheetMode int

const (
	NORMAL_MODE SheetMode = iota
	INSERT_MODE SheetMode = iota
	EXIT_MODE   SheetMode = iota
	YANK_MODE   SheetMode = iota
	PUT_MODE    SheetMode = iota
)

func processTermboxEvents(s *sheet.Sheet) {
	prompt := ""
	stringEntry := false
	smode := NORMAL_MODE
	valBuffer := bytes.Buffer{}
	insAlign := align.AlignRight

	// Display
	go func() {
		for _ = range time.Tick(200 * time.Millisecond) {
			switch smode {
			case NORMAL_MODE:
				selSel, _ := s.GetCell(s.SelectedCell)
				display.DisplayValue(fmt.Sprintf("%s (%s) [%s]", s.SelectedCell, s.DisplayFormat(s.SelectedCell), selSel.StatusBarVal()), 0, 0, 80, align.AlignLeft, false)
			case INSERT_MODE:
				display.DisplayValue(fmt.Sprintf("i> %s %s = %s", prompt, s.SelectedCell, valBuffer.String()), 0, 0, 80, align.AlignLeft, false)
			case EXIT_MODE:
				display.DisplayValue(fmt.Sprintf("File \"%s\" is modified, save before exiting?", s.Filename), 0, 0, 80, align.AlignLeft, false)
			case YANK_MODE:
				display.DisplayValue("Yank row/column:  r: row  c: column", 0, 0, 80, align.AlignLeft, false)
			case PUT_MODE:
				display.DisplayValue("Put row/column:  r: row  c: column", 0, 0, 80, align.AlignLeft, false)
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
						insAlign = align.AlignRight
					case '<':
						prompt = "leftstring"
						smode = INSERT_MODE
						insAlign = align.AlignLeft
						stringEntry = true
					case '>':
						prompt = "rightstring"
						smode = INSERT_MODE
						insAlign = align.AlignRight
						stringEntry = true
					case '\\':
						prompt = "label"
						smode = INSERT_MODE
						insAlign = align.AlignCenter
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
						s.ClearCell(s.SelectedCell)
					case 'y':
						smode = YANK_MODE
					case 'p':
						smode = PUT_MODE
					}
				}
			case INSERT_MODE:
				if ev.Key == termbox.KeyEnter {
					s.SetCell(s.SelectedCell, sheet.NewCell(valBuffer.String(), insAlign, stringEntry))
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
					s.Save()
				}
				termbox.Close()
				return
			case YANK_MODE:
				if ev.Key == 0 && ev.Ch == 'r' {
					s.YankRow()
				} else if ev.Key == 0 && ev.Ch == 'c' {
					s.YankColumn()
				}
				smode = NORMAL_MODE

			case PUT_MODE:
				if ev.Key == 0 && ev.Ch == 'r' {
					s.PutRow()
				} else if ev.Key == 0 && ev.Ch == 'c' {
					s.PutColumn()
				}
				smode = NORMAL_MODE
			}
		}
	}
}
