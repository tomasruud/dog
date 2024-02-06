package ansi

import (
	"image"
	"image/color"
	"io"
	"strings"

	ctools "github.com/gookit/color"
	"golang.org/x/image/draw"
)

// Encoder configures encoding for images in ansi escape codes.
type Encoder struct {
	MaxW       int
	MaxH       int
	ColorLevel ctools.Level
	Matte      color.Color
	Scaler     draw.Scaler
}

func (e Encoder) scaler() draw.Scaler {
	if e.Scaler == nil {
		return draw.BiLinear
	}
	return e.Scaler
}

func (e Encoder) matte() color.Color {
	if e.Matte == nil {
		return color.White
	}
	return e.Matte
}

func (e Encoder) Encode(w io.Writer, m image.Image) error {
	size := m.Bounds().Size()
	for e.MaxW < size.X || e.MaxH < (size.Y/2) {
		size.X = size.X / 2
		size.Y = size.Y / 2
	}

	dst := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))

	e.scaler().Scale(dst, dst.Rect, m, m.Bounds(), draw.Over, nil)

	buf := strings.Builder{}

	bounds := dst.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			bg := e.escape(dst.At(x, y), true)
			buf.WriteString(ctools.StartSet)
			buf.WriteString(bg)
			buf.WriteString("m")

			fg := e.escape(dst.At(x, y+1), false)
			buf.WriteString(ctools.StartSet)
			buf.WriteString(fg)
			buf.WriteString("m")

			buf.WriteString("▄")
		}

		buf.WriteString(ctools.ResetSet)
		buf.WriteString("\n")
	}

	_, err := w.Write([]byte(buf.String()))
	return err
}

// escape creates ansi escape code, fitting the color to the current color level.
func (e Encoder) escape(c color.Color, bg bool) string {
	rgb := ctools.RGB(e.rgb(c))

	if bg {
		rgb = rgb.ToBg()
	} else {
		rgb = rgb.ToFg()
	}

	switch e.ColorLevel {
	case ctools.Level16:
		return rgb.C16().String()

	case ctools.Level256:
		return rgb.C256().String()
	}

	return rgb.String()
}

// rgb compiles rgba colors on a predefined background color.
// Based on the algorithm suggested here https://stackoverflow.com/a/746937/738608
func (e Encoder) rgb(c color.Color) (r, g, b uint8) {
	or, og, ob, oa := c.RGBA()
	mr, mg, mb, _ := e.matte().RGBA()

	alpha := float32(oa>>8) / 255

	r = uint8(alpha*float32(or>>8) + (1-alpha)*float32(mr>>8))
	g = uint8(alpha*float32(og>>8) + (1-alpha)*float32(mg>>8))
	b = uint8(alpha*float32(ob>>8) + (1-alpha)*float32(mb>>8))

	return r, g, b
}
