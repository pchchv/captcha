package captcha

import "math"

type hsva struct {
	h float64
	s float64
	v float64
	a uint8
}

func (c hsva) RGBA() (r, g, b, a uint32) {
	var (
		i                = math.Floor(c.h * 6)
		f                = c.h*6 - i
		p                = c.v * (1.0 - c.s)
		q                = c.v * (1.0 - f*c.s)
		t                = c.v * (1 - (1-f)*c.s)
		red, green, blue float64
	)
	switch int(i) % 6 {
	case 0:
		red, green, blue = c.v, t, p
	case 1:
		red, green, blue = q, c.v, p
	case 2:
		red, green, blue = p, c.v, t
	case 3:
		red, green, blue = p, q, c.v
	case 4:
		red, green, blue = t, p, c.v
	case 5:
		red, green, blue = c.v, p, q
	}

	r = uint32(red * 255)
	r |= r << 8
	g = uint32(green * 255)
	g |= g << 8
	b = uint32(blue * 255)
	b |= b << 8
	a = uint32(c.a)
	a |= a << 8

	return
}
