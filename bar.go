package giocharts

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Bar struct {
	Data  []float64
	Label func(i int) string
	Color func(i int) color.NRGBA
	Theme *material.Theme
}

func (b Bar) Layout(gtx layout.Context) layout.Dimensions {
	if b.Label == nil {
		b.Label = func(i int) string {
			return fmt.Sprint(i)
		}
	}
	if b.Color == nil {
		b.Color = func(i int) color.NRGBA {
			return hslToRGB(float64(i)/float64(len(b.Data)), .8, .8)
		}
	}
	if b.Theme == nil {
		b.Theme = defaultTheme()
	}

	b.renderBars(gtx)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func (b Bar) renderBars(gtx layout.Context) {
	n := len(b.Data)
	boxSize := gtx.Constraints.Min.X / n

	_, dataMax := minMax(b.Data)
	dmx := dataMax
	// dmi := dataMin
	dataRange := dmx

	maxBoxHeight := (dmx / dataRange) * float64(gtx.Constraints.Max.Y)

	for i, n := range b.Data {
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
		if i == len(b.Data)-1 {
			bar.Right = unit.Value{V: 0, U: unit.UnitDp}
		}

		bar.Layout(gtx, fillWithLabel(b.Theme, b.Label(i), b.Color(i)))

		trans.Pop()
	}
}
