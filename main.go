package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/scrouthtv/gosc/internal/sheet"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		fmt.Println("Could not start termbox.")
		os.Exit(1)
	}

	flag.Parse()

	// For now, only argument is sheet file name.
	sheet := sheet.NewSheet(flag.Arg(0))

	processTermboxEvents(&sheet)
}
