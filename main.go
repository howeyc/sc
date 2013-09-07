package main

import (
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
