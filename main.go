package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/nsf/termbox-go"
	"gopkg.in/yaml.v2"
)

var script = `Welcome to the thing!`

type pres struct {
	Header string  `yaml:"header"`
	Footer string  `yaml:"footer"`
	Slides []slide `yaml:"slides"`
}

type slide struct {
	Header  string `yaml:"header"`
	Footer  string `yaml:"footer"`
	Content string `yaml:"content"`
}

var fgCols = []termbox.Attribute{
	termbox.ColorRed,
	termbox.ColorGreen,
	termbox.ColorYellow,
	termbox.ColorBlue,
	termbox.ColorMagenta,
	termbox.ColorCyan,
	termbox.ColorWhite,
}

var bgCols = []termbox.Attribute{
	termbox.ColorDefault,
	termbox.ColorBlack,
	termbox.ColorWhite,
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	defer termbox.Close()

	f, err := os.Open("slides.yaml")
	if err != nil {
		panic(err)
	}

	p := pres{}
	d := yaml.NewDecoder(f)
	err = d.Decode(&p)
	if err != nil {
		panic(err)
	}

	s := newScreen(termbox.ColorWhite, termbox.ColorBlack)

	for i := 0; i < len(p.Slides); i++ {
		slide := p.Slides[i]

		next := newScreen(fgCols[rand.Intn(len(fgCols))], bgCols[rand.Intn(len(bgCols))])
		next.border(fgCols[rand.Intn(len(fgCols))])

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

		event := termbox.PollEvent()

		if event.Type == termbox.EventKey {
			if event.Ch == 'q' {
				return
			} else if event.Key == termbox.KeyArrowLeft {
				i -= 2
			}
		}

		if i == len(p.Slides)-1 {
			i = -1
		}
	}
}
