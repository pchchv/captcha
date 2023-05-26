// Package provides an simple, unopinionated API for captcha generation
package captcha

import "image/color"

const charPreset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Options manage captcha generation.
type Options struct {
	// BackgroundColor is captcha image's background color.
	// By default it is color.Transparent.
	BackgroundColor color.Color
	// CharPreset defines the text on the captcha image.
	// By default:
	// ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789
	CharPreset string
	// TextLength is the length of the captcha text.
	// By default it is 4.
	TextLength int
	// CurveNumber is the number of curves to draw on captcha image.
	// By default it is 2.
	CurveNumber int
	// FontDPI controls DPI (dots per inch) of font.
	// By default it is 72.0.
	FontDPI float64
	// FontScale controls the scale of font.
	// By default it is 1.0.
	FontScale float64
	// Noise controls the amount of noise to be drawn.
	// By default, a noise point is drawn every 28 pixels.
	// By default it is 1.0.
	Noise float64
	// Palette is the set of colors to chose from.
	Palette color.Palette
	width   int
	height  int
}

func newDefaultOption(width, height int) *Options {
	return &Options{
		BackgroundColor: color.Transparent,
		CharPreset:      charPreset,
		TextLength:      4,
		CurveNumber:     2,
		FontDPI:         72.0,
		FontScale:       1.0,
		Noise:           1.0,
		Palette:         []color.Color{},
		width:           width,
		height:          height,
	}
}
