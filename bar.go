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
