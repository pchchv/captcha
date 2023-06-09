// Package provides an simple, unopinionated API for captcha generation
package captcha

import (
	_ "embed"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const charPreset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var (
	//go:embed fonts/Comismsh.ttf
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

// New creates a new captcha.
func New(width int, height int, option ...SetOption) (*Captcha, error) {
	options := newDefaultOption(width, height)

	for _, setOption := range option {
		setOption(options)
	}

	text := randomText(options)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	if err := drawWithOption(text, img, options); err != nil {
		return nil, err
	}

	return &Captcha{Text: text, img: img}, nil
}

// NewMathExpr creates a new captcha.
// Generates an image with a mathematical expression like `1 + 2`.
func NewMathExpr(width int, height int, option ...SetOption) (*Captcha, error) {
	options := newDefaultOption(width, height)

	for _, setOption := range option {
		setOption(options)
	}

	text, equation := randomEquation()
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	if err := drawWithOption(equation, img, options); err != nil {
		return nil, err
	}

	return &Captcha{Text: text, img: img}, nil
}

// NewCustomGenerator creates a new captcha based on a custom text generator.
func NewCustomGenerator(width int, height int, generator func() (anwser string, question string), option ...SetOption) (*Captcha, error) {
	options := newDefaultOption(width, height)

	for _, setOption := range option {
		setOption(options)
	}

	answer, question := generator()
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	if err := drawWithOption(question, img, options); err != nil {
		return nil, err
	}

	return &Captcha{Text: answer, img: img}, nil
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

	return text
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

	return text, equation
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

	return (float64(max) + float64(min)) / (2 * 255)
}

func drawWithOption(text string, img *image.NRGBA, options *Options) error {
	draw.Draw(img, img.Bounds(), &image.Uniform{options.BackgroundColor}, image.Point{}, draw.Src)
	drawNoise(img, options)
	drawCurves(img, options)

	return drawText(text, img, options)
}

func drawCurves(img *image.NRGBA, opts *Options) {
	for i := 0; i < opts.CurveNumber; i++ {
		drawSineCurve(img, opts)
	}
}

func drawSineCurve(img *image.NRGBA, opts *Options) {
	var xStart, xEnd int
	if opts.width <= 40 {
		xStart, xEnd = 1, opts.width-1
	} else {
		xStart = rand.Intn(opts.width/10) + 1
		xEnd = opts.width - rand.Intn(opts.width/10) - 1
	}

	curveHeight := float64(rand.Intn(opts.height/6) + opts.height/6)
	yStart := rand.Intn(opts.height*2/3) + opts.height/6
	angle := 1.0 + rand.Float64()
	yFlip := 1.0
	if rand.Intn(2) == 0 {
		yFlip = -1.0
	}

	curveColor := randomColorFromOptions(opts)

	for x1 := xStart; x1 <= xEnd; x1++ {
		y := math.Sin(math.Pi*angle*float64(x1)/float64(opts.width)) * curveHeight * yFlip
		img.Set(x1, int(y)+yStart, curveColor)
	}
}

func drawText(text string, img *image.NRGBA, opts *Options) error {
	ctx := freetype.NewContext()
	ctx.SetDPI(opts.FontDPI)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)
	ctx.SetHinting(font.HintingFull)
	ctx.SetFont(ttfFont)

	fontSpacing := opts.width / len(text)
	fontOffset := rand.Intn(fontSpacing / 2)

	for idx, char := range text {
		fontScale := 0.8 + rand.Float64()*0.4
		fontSize := float64(opts.height) / fontScale * opts.FontScale
		ctx.SetFontSize(fontSize)
		ctx.SetSrc(image.NewUniform(randomColorFromOptions(opts)))
		x := fontSpacing*idx + fontOffset
		y := opts.height/6 + rand.Intn(opts.height/3) + int(fontSize/2)
		pt := freetype.Pt(x, y)
		if _, err := ctx.DrawString(string(char), pt); err != nil {
			return err
		}
	}

	return nil
}

func drawNoise(img *image.NRGBA, opts *Options) {
	noiseCount := (opts.width * opts.height) / int(28.0/opts.Noise)

	for i := 0; i < noiseCount; i++ {
		x := rand.Intn(opts.width)
		y := rand.Intn(opts.height)
		img.Set(x, y, randomColor())
	}
}
