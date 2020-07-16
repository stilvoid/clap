package pres

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Slide struct {
	Header  string `yaml:"header"`
	Footer  string `yaml:"footer"`
	Fg      string `yaml:"fg"`
	Bg      string `yaml:"bg"`
	Content string `yaml:"content"`
	Theme   string `yaml:"theme"`
}

type Pres struct {
	Header string  `yaml:"header"`
	Footer string  `yaml:"footer"`
	Fg     string  `yaml:"fg"`
	Bg     string  `yaml:"bg"`
	Slides []Slide `yaml:"slides"`
}

func (p Pres) Run() {
	defer termbox.Close()

	setup()

	s := newScreen(termbox.ColorWhite, termbox.ColorBlack)

	/*
		s.words("Loading...")
		for i := 0; i < 32; i++ {
			s.border(termbox.Attribute(rand.Intn(16)))
			s.display()
			time.Sleep(time.Second / 16)
		}
	*/

	for i := 0; i < len(p.Slides); i++ {
		slide := p.Slides[i]

		fgCol, bgCol := s.fg, s.bg

		// Parse colours
		fg, bg := slide.Fg, slide.Bg
		if fg == "" {
			fg = p.Fg
		}
		if bg == "" {
			bg = p.Bg
		}

		if fg != "" {
			fgCol = colours[fg]
		}

		if bg != "" {
			bgCol = colours[bg]
		}

		next := newScreen(fgCol, bgCol)
		next.border(fgCol)

		h := slide.Header
		if h == "" {
			h = p.Header
		}
		if h != "" {
			next.header(h)
		}

		f := slide.Footer
		if f == "" {
			f = p.Footer
		}
		if f != "" {
			next.footer(f)
		}

		next.words(slide.Content)

		next.page(i+1, len(p.Slides))

		switch rand.Intn(4) {
		case 0:
			s.s(next)
		case 1:
			s.se(next)
		case 2:
			s.n(next)
		case 3:
			s.sw(next)
		}

		if slide.Theme == "rain" {
			s.rain()
		} else if slide.Theme == "rainbow" {
			s.rainbow()
		}

		event := termbox.PollEvent()

		if event.Type == termbox.EventKey {
			if event.Ch == 'q' {
				return
			} else if event.Key == termbox.KeyArrowLeft {
				i = i - 2
				if i < -1 {
					i = len(p.Slides) + i
				}
			}
		}

		if i == len(p.Slides)-1 {
			i = -1
		}
	}
}
