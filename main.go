package main

import (
	"flag"
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

	flag.Parse()

	sheet := newSheet(flag.Arg(0))

	sheet.display()

	processTermboxEvents(&sheet)
}
