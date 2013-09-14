Spreadsheet Calculator
========

This is a text-based spreadsheet program compatible with [sc](http://www.ibiblio.org/pub/Linux/apps/financial/spreadsheet/sc-7.16.lsm) files. sc is a public domain curses-based speadsheet calculator program that has similar key bindings to vi and works on Unix-like OSs. This program supports the features of sc I use, can read/save sc files and works anywhere [termbox-go](https://github.com/nsf/termbox-go) is supported.

Essentially this allows me to view/edit sc files on a Windows machine.

## Supported Features

* Vi-like key movements
* Expressions (B23+E23*5)

## Normal Mode Commands

* Movement: h, j, k, l
* Enter insert mode for numbers and expressions: =, i
* Enter insert mode for strings with alignment: < (Left-Aligned), > (Right-Aligned), \ (Centered)
* Delete (clear) cell: x
* Quit: q

## TODO

* Commands: (y)ank, (p)ut
* Functions for expressions.