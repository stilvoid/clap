package pres

import "github.com/nsf/termbox-go"

var colours = map[string]termbox.Attribute{
	"grey":    termbox.ColorDefault,
	"black":   termbox.ColorBlack,
	"red":     termbox.ColorRed,
	"green":   termbox.ColorGreen,
	"yellow":  termbox.ColorYellow,
	"blue":    termbox.ColorBlue,
	"magenta": termbox.ColorMagenta,
	"cyan":    termbox.ColorCyan,
	"white":   termbox.ColorWhite,
}

var rainbow = []termbox.Attribute{
	termbox.ColorRed | termbox.AttrBold,
	termbox.ColorYellow | termbox.AttrBold,
	termbox.ColorMagenta,
	termbox.ColorGreen,
	termbox.ColorCyan | termbox.AttrBold,
	termbox.ColorYellow,
	termbox.ColorBlue | termbox.AttrBold,
}
