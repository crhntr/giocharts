package giocharts

import (
	"image"
	"image/color"
	"math"

	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"golang.org/x/exp/constraints"
)

var black = hslToRGB(0, 0, 0)

func defaultTheme() *material.Theme {
	return material.NewTheme(gofont.Collection())
}

func fillWithLabel(th *material.Theme, content string, backgroundColor color.NRGBA) func(gtx layout.Context) layout.Dimensions {
	return func(gtx layout.Context) layout.Dimensions {
		colorBox(gtx, gtx.Constraints.Max, backgroundColor)
		body := material.Body1(th, content)
		body.Alignment = text.Middle
		return layout.UniformInset(unit.Px(5)).Layout(gtx, body.Layout)
	}
}

func colorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

func minMax[N constraints.Ordered](numbers []N) (min, max N) {
	if len(numbers) == 0 {
		return
	}
	min = numbers[0]
	max = numbers[0]
	for _, n := range numbers[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	return
}

func hslToRGB(h, s, l float64) color.NRGBA {
	fixHue := func(p, q, t float64) float64 {
		if t < 0 {
			t += 1
		}
		if t > 1 {
			t -= 1
		}
		if t < 1.0/6 {
			return p + (q-p)*6*t
		}
		if t < 1.0/2 {
			return q
		}
		if t < 2.0/3 {
			return p + (q-p)*(2.0/3-t)*6
		}
		return p
	}

	rgbFloatToInt8 := func(r, g, b float64) color.NRGBA {
		col := color.NRGBA{
			R: uint8(math.Round(r * 255)),
			G: uint8(math.Round(g * 255)),
			B: uint8(math.Round(b * 255)),
			A: 255,
		}
		return col
	}

	if s == 0 {
		return rgbFloatToInt8(l, l, l)
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}

	p := 2*l - q
	r := fixHue(p, q, h+1.0/3)
	g := fixHue(p, q, h)
	b := fixHue(p, q, h-1.0/3)

	return rgbFloatToInt8(r, g, b)
}
