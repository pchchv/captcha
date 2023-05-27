package captcha

import (
	"image/color"
	"image/color/palette"
	"testing"
)

func TestNewCaptchaOptions(t *testing.T) {
	New(100, 34, func(options *Options) {
		options.BackgroundColor = color.Opaque
		options.CharPreset = "1234567890"
		options.CurveNumber = 0
		options.TextLength = 6
		options.Palette = palette.WebSafe
	})

	NewMathExpr(100, 34, func(options *Options) {
		options.BackgroundColor = color.Black
	})

	NewCustomGenerator(100, 34, func() (anwser string, question string) {
		return "4", "2x2?"
	}, func(o *Options) {
		o.BackgroundColor = color.Black
	})
}

func TestNilFontError(t *testing.T) {
	temp := ttfFont
	ttfFont = nil

	_, err := New(150, 50)
	if err == nil {
		t.Fatal("Expect to get nil font error")
	}

	_, err = NewMathExpr(150, 50)
	if err == nil {
		t.Fatal("Expect to get nil font error")
	}

	_, err = NewCustomGenerator(150, 50, func() (anwser string, question string) {
		return "1", "2"
	})
	if err == nil {
		t.Fatal("Expect to get nil font error")
	}

	ttfFont = temp
}
