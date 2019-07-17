package main

import (
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

var w, h int

var wipeTime = time.Second.Nanoseconds() / 2

func init() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	w, h = termbox.Size()
}

type screen struct {
	cells []termbox.Cell
	fg    termbox.Attribute
	bg    termbox.Attribute
}

func newScreen(fg, bg termbox.Attribute) screen {
	s := screen{
		cells: make([]termbox.Cell, w*h),
		fg:    fg,
		bg:    bg,
	}

	for i := 0; i < w*h; i++ {
		s.cells[i] = termbox.Cell{' ', fg, bg}
	}

	return s
}

func (s screen) display() {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := s.cells[y*w+x]
			termbox.SetCell(x, y, c.Ch, c.Fg, c.Bg)
		}
	}

	termbox.Flush()
}

func (s screen) border(a termbox.Attribute) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if y < 1 || x <= 1 || y > h-2 || x >= w-2 {
				s.cells[y*w+x] = termbox.Cell{'█', a, termbox.ColorDefault}
			}
		}
	}
}

func (s screen) line(y int, line string, a termbox.Attribute) {
	ox := (w - len(line)) / 2

	for x, c := range line {
		s.cells[y*w+x+ox] = termbox.Cell{c, a, s.bg}
	}
}

func (s screen) words(words string) {
	lines := strings.Split(words, "\n")

	oy := (h - len(lines)) / 2

	for y, line := range lines {
		a := s.fg
		if y == 0 {
			a |= termbox.AttrBold
		}

		s.line(y+oy, line, a)
	}
}

func (s screen) header(header string) {
	s.line(2, header, s.fg)
}

func (s screen) footer(footer string) {
	s.line(h-3, footer, s.fg)
}

func (s screen) s(next screen) {
	t := time.Duration(wipeTime / int64(h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := next.cells[y*w+x]
			s.cells[y*w+x] = c
			termbox.SetCell(x, y, c.Ch, c.Fg, c.Bg)
		}
		termbox.Flush()
		time.Sleep(t)
	}
}

func (s screen) n(next screen) {
	t := time.Duration(wipeTime / int64(h))

	for y := h - 1; y >= 0; y-- {
		for x := 0; x < w; x++ {
			c := next.cells[y*w+x]
			s.cells[y*w+x] = c
			termbox.SetCell(x, y, c.Ch, c.Fg, c.Bg)
		}
		termbox.Flush()
		time.Sleep(t)
	}
}

func (s screen) se(next screen) {
	t := time.Duration(wipeTime / int64(w+h))

	for i := 0; i < w+h; i++ {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if x+y == i {
					c := next.cells[y*w+x]
					s.cells[y*w+x] = c
					termbox.SetCell(x, y, c.Ch, c.Fg, c.Bg)
				}
			}
		}
		termbox.Flush()
		time.Sleep(t)
	}
}

func (s screen) sw(next screen) {
	t := time.Duration(wipeTime / int64(w+h))

	for i := 0; i < w+h; i++ {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				if (w-x)+y == i {
					c := next.cells[y*w+x]
					s.cells[y*w+x] = c
					termbox.SetCell(x, y, c.Ch, c.Fg, c.Bg)
				}
			}
		}
		termbox.Flush()
		time.Sleep(t)
	}
}
