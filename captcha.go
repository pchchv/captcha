// Package provides an simple, unopinionated API for captcha generation
package captcha

import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

const charPreset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var (
	ttf     []byte
	ttfFont *truetype.Font
)

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

// Captcha is the result of captcha generation.
// It has a `Text` field and a private `img` field,
// which will be used in the `WriteImage` receiver.
type Captcha struct {
	Text string
	img  *image.NRGBA
}

// SetOption is a function used to change the default settings.
type SetOption func(*Options)

func init() {
	ttfFont, _ = freetype.ParseFont(ttf)
	rand.Seed(time.Now().UnixNano())
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

// WriteImage encodes image data and writes it to io.Writer.
// Returns an error possible when encoding PNG.
func (c *Captcha) WriteImage(w io.Writer) error {
	return png.Encode(w, c.img)
}

// WriteJPG encodes the image data into JPEG format and writes it to io.Writer.
// Returns an error possible when encoding JPEG.
func (c *Captcha) WriteJPG(w io.Writer, o *jpeg.Options) error {
	return jpeg.Encode(w, c.img, o)
}

// WriteGIF encodes the image data into GIF format and writes it to io.Writer.
// Returns an error possible when encoding GIF.
func (c *Captcha) WriteGIF(w io.Writer, o *gif.Options) error {
	return gif.Encode(w, c.img, o)
}

// LoadFont lets load an external font.
func LoadFont(fontData []byte) (err error) {
	ttfFont, err = freetype.ParseFont(fontData)
	return err
}

// LoadFontFromReader load an external font from an io.Reader interface.
func LoadFontFromReader(reader io.Reader) error {
	b, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	return LoadFont(b)
}

func randomText(opts *Options) (text string) {
	n := len([]rune(opts.CharPreset))
	for i := 0; i < opts.TextLength; i++ {
		text += string([]rune(opts.CharPreset)[rand.Intn(n)])
	}

	return
}

func randomColor() color.RGBA {
	red := rand.Intn(256)
	green := rand.Intn(256)
	blue := rand.Intn(256)

	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: uint8(255)}
}

func randomColorFromOptions(opts *Options) color.Color {
	length := len(opts.Palette)
	if length == 0 {
		return randomInvertColor(opts.BackgroundColor)
	}

	return opts.Palette[rand.Intn(length)]
}

func randomEquation() (text string, equation string) {
	left := 1 + rand.Intn(9)
	right := 1 + rand.Intn(9)
	text = strconv.Itoa(left + right)
	equation = strconv.Itoa(left) + "+" + strconv.Itoa(right)

	return
}

func randomInvertColor(base color.Color) color.Color {
	var value float64
	baseLightness := getLightness(base)
	if baseLightness >= 0.5 {
		value = baseLightness - 0.3 - rand.Float64()*0.2
	} else {
		value = baseLightness + 0.3 + rand.Float64()*0.2
	}
	hue := float64(rand.Intn(361)) / 360
	saturation := 0.6 + rand.Float64()*0.2

	return hsva{h: hue, s: saturation, v: value, a: 255}
}

func maxColor(numList ...uint32) (max uint32) {
	for _, num := range numList {
		colorVal := num & 255
		if colorVal > max {
			max = colorVal
		}
	}

	return max
}

func minColor(numList ...uint32) (min uint32) {
	min = 255
	for _, num := range numList {
		colorVal := num & 255
		if colorVal < min {
			min = colorVal
		}
	}

	return min
}

func getLightness(colour color.Color) float64 {
	r, g, b, a := colour.RGBA()
	if a == 0 {
		return 1.0
	}

	max := maxColor(r, g, b)
	min := minColor(r, g, b)

	l := (float64(max) + float64(min)) / (2 * 255)

	return l
}
