package captcha

import (
	"errors"
	"image/color"
	"image/color/palette"
	"os"
	"testing"

	"golang.org/x/image/font/gofont/goregular"
)

type errReader struct{}

func (errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("")
}

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

func TestReaderErr(t *testing.T) {
	if err := LoadFontFromReader(errReader{}); err == nil {
		t.Fatal("Expect to get io.Reader error")
	}
}

func TestLoadFont(t *testing.T) {
	if err := LoadFont(goregular.TTF); err != nil {
		t.Fatal("Fail to load go font")
	}

	if err := LoadFont([]byte("invalid")); err == nil {
		t.Fatal("LoadFont incorrectly parse an invalid font")
	}
}

func TestLoadFontFromReader(t *testing.T) {
	file, err := os.Open("./fonts/Comismsh.ttf")
	if err != nil {
		t.Fatal("Fail to load test file")
	}

	if err = LoadFontFromReader(file); err != nil {
		t.Fatal("Fail to load font from io.Reader")
	}
}
