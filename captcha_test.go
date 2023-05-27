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
