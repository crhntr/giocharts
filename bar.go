package giocharts

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"golang.org/x/exp/constraints"
)

type Bar struct {
	Data  []float64
	Label func(i int) string
	Color func(i int) color.NRGBA
}

func (s Bar) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	n := len(s.Data)
	boxSize := gtx.Constraints.Min.X / n

	_, dataMax := minMax(s.Data)
	dmx := dataMax
	// dmi := dataMin
	dataRange := dmx

	maxBoxHeight := (dmx / dataRange) * float64(gtx.Constraints.Max.Y)

	labelFunc := s.Label
	if labelFunc == nil {
		labelFunc = func(i int) string {
			return fmt.Sprint(i)
		}
	}
	colorFunc := s.Color
	if colorFunc == nil {
		colorFunc = func(i int) color.NRGBA {
			return hslToRGB(float64(i)/float64(len(s.Data)), .8, .8)
		}
	}

	for i, n := range s.Data {
		gtx := gtx
		boxHeight := (n / dataRange) * float64(gtx.Constraints.Max.Y)
		gtx.Constraints = layout.Exact(image.Pt(boxSize, int(boxHeight)))
		trans := op.Offset(f32.Pt(float32(boxSize*i), float32(maxBoxHeight-boxHeight))).Push(gtx.Ops)

		bar := layout.Inset{
			Top:   unit.Value{V: 5, U: unit.UnitDp},
			Left:  unit.Value{V: 5, U: unit.UnitDp},
			Right: unit.Value{V: 5, U: unit.UnitDp},
		}

		if i == 0 {
			bar.Left = unit.Value{V: 0, U: unit.UnitDp}
		}
		if i == len(s.Data)-1 {
			bar.Right = unit.Value{V: 0, U: unit.UnitDp}
		}

		bar.Layout(gtx, fillWithLabel(th, labelFunc(i), colorFunc(i)))

		trans.Pop()
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func fillWithLabel(th *material.Theme, text string, backgroundColor color.NRGBA) func(gtx layout.Context) layout.Dimensions {
	return func(gtx layout.Context) layout.Dimensions {
		colorBox(gtx, gtx.Constraints.Max, backgroundColor)
		return layout.Inset{
			Top:    unit.Value{V: 20, U: unit.UnitDp},
			Left:   unit.Value{V: 20, U: unit.UnitDp},
			Right:  unit.Value{V: 20, U: unit.UnitDp},
			Bottom: unit.Value{V: 20, U: unit.UnitDp},
		}.Layout(gtx, material.Body1(th, text).Layout)
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
