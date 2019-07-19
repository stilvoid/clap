package main

import (
	"fmt"
	"math"
	"math/rand"
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

	termbox.SetOutputMode(termbox.Output256)
}

type screen struct {
	cells     []termbox.Cell
	fg        termbox.Attribute
	bg        termbox.Attribute
	animating bool
	stopped   bool
}

func newScreen(fg, bg termbox.Attribute) screen {
	s := screen{
		cells:     make([]termbox.Cell, w*h),
		fg:        fg,
		bg:        bg,
		animating: false,
		stopped:   true,
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
	c := termbox.Cell{'█', a, termbox.ColorDefault}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if y < 1 || x <= 1 || y > h-2 || x >= w-2 {
				s.cells[y*w+x] = c
			}
		}
	}
}

func (s screen) line(y int, line string, a termbox.Attribute, maxLength int) {
	ox := (w - len(line)) / 2

	if strings.HasPrefix(line, "  ") {
		ox = (w - maxLength) / 2
	}

	for x, c := range line {
		s.cells[y*w+x+ox] = termbox.Cell{c, a, s.bg}
	}
}

func (s screen) words(words string) {
	lines := strings.Split(words, "\n")

	oy := (h - len(lines)) / 2

	maxLength := 0
	for _, line := range lines {
		maxLength = int(math.Max(float64(maxLength), float64(len(line))))
	}

	for y, line := range lines {
		a := s.fg

		if y == 0 {
			a |= termbox.AttrBold
		}

		s.line(y+oy, line, a, maxLength)
	}
}

func (s screen) page(n, m int) {
	display := fmt.Sprintf("%d / %d", n, m)

	ox := w - len(display) - 4

	for x, c := range display {
		s.cells[(h-3)*w+ox+x] = termbox.Cell{c, s.fg, s.bg}
	}
}

func (s screen) header(header string) {
	s.line(2, header, s.fg, len(header))
}

func (s screen) footer(footer string) {
	s.line(h-3, footer, s.fg, len(footer))
}

func (s *screen) replace(next screen) {
	s.cells = next.cells
	s.fg = next.fg
	s.bg = next.bg
	s.animating = false

	for !s.stopped {
		time.Sleep(time.Second / 10)
	}
	s.display()
	termbox.Sync()
}

func (s *screen) s(next screen) {
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

	s.replace(next)
}

func (s *screen) n(next screen) {
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

	s.replace(next)
}

func (s *screen) se(next screen) {
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

	s.replace(next)
}

func (s *screen) sw(next screen) {
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

	s.replace(next)
}

func (s *screen) rain() {
	s.animating = true
	s.stopped = false

	go func() {
		type drop struct {
			x, y int
			dead bool
		}

		drops := make([]*drop, 0)

		for s.animating {
			if rand.Float64() < 0.3 {
				drops = append(drops, &drop{2 + rand.Intn(w-4), 1, false})
			}

			for _, drop := range drops {
				if drop.dead {
					continue
				}

				next := s.cells[drop.y*w+drop.x]
				termbox.SetCell(drop.x, drop.y, '\'', termbox.ColorBlue, next.Bg)
			}

			termbox.Flush()

			time.Sleep(time.Second / 20)

			for _, drop := range drops {
				if drop.dead {
					continue
				}

				previous := s.cells[drop.y*w+drop.x]
				termbox.SetCell(drop.x, drop.y, previous.Ch, previous.Fg, previous.Bg)

				drop.y++

				if drop.y >= h-1 {
					drop.dead = true
				}
			}
		}

		s.stopped = true
	}()
}

func (s *screen) rainbow() {
	s.animating = true
	s.stopped = false

	go func() {
		count := 0

		for s.animating {
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					cell := s.cells[y*w+x]
					termbox.SetCell(x, y, cell.Ch, rainbow[(y+count)%len(rainbow)], cell.Bg)
				}
			}

			termbox.Flush()

			time.Sleep(time.Second / 2)

			count++
		}

		s.stopped = true
	}()
}
