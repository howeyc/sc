package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		fmt.Println("Could not start termbox.")
		os.Exit(1)
	}

	sheet := Sheet{}
	sheet.display()
	time.Sleep(time.Second * 10)
	termbox.Close()
}
