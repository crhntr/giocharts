package giocharts

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Bar struct {
	Data         []float64
	BoxLabel     func(i int) string
	BarColor     func(i int) color.NRGBA
	SpaceBetween float32
	Theme        *material.Theme
}

func (b Bar) Layout(gtx layout.Context) layout.Dimensions {
	if b.BoxLabel == nil {
		b.BoxLabel = func(i int) string {
			if len(b.Data)*200 > gtx.Constraints.Max.X {
				return ""
			}
			return fmt.Sprintf("%2.2f", b.Data[i])
		}
	}
	if b.BarColor == nil {
		b.BarColor = func(i int) color.NRGBA {
			return hslToRGB(float64(i)/float64(len(b.Data)), .8, .8)
		}
	}
	if b.Theme == nil {
		b.Theme = defaultTheme()
	}
	if b.SpaceBetween == 0 {
		b.SpaceBetween = 10
	}

	dMinValue, dMaxValue := minMax(b.Data)

	return b.renderYLegend(gtx, b.Theme, dMinValue, dMaxValue)
}

func (b Bar) renderYLegend(gtx layout.Context, th *material.Theme, minValue, maxValue float64) layout.Dimensions {
	offset := unit.Dp(50)
	prefixWidth := gtx.Px(offset)
	// lineWidth := gtx.Px(unit.Dp(5))

	maxX := gtx.Constraints.Max.X
	gtx.Constraints.Max.X = prefixWidth
	writeYLegendLine(gtx)
	gtx.Constraints.Max.X = maxX

	//labelFunc := func(v float64) layout.Widget {
	//	l := material.Label(th, th.TextSize, fmt.Sprintf("%.2f", maxValue))
	//	l.Alignment = text.End
	//	return l.Layout
	//}

	//correction := (maxValue - minValue) / maxValue * float64(gtx.Constraints.Max.Y)

	_ = layout.Inset{
		Top:    th.TextSize,
		Bottom: th.TextSize,
		Left:   offset,
	}.Layout(gtx, b.renderBars(maxValue))

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func writeYLegendLine(gtx layout.Context) layout.Dimensions {
	size := image.Pt(2, gtx.Constraints.Max.Y)
	defer clip.Rect{
		Min: image.Pt(gtx.Constraints.Max.X-10, 0),
		Max: image.Pt(gtx.Constraints.Max.X, gtx.Constraints.Max.Y),
	}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: black}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

func (b Bar) renderBars(maxN float64) func(gtx layout.Context) layout.Dimensions {
	return func(gtx layout.Context) layout.Dimensions {
		n := len(b.Data)
		boxWidth := gtx.Constraints.Max.X / n

		for i, n := range b.Data {
			gtx := gtx
			boxHeight := (n / maxN) * float64(gtx.Constraints.Max.Y)
			offsetTop := float32(gtx.Constraints.Max.Y) - float32(boxHeight)
			offsetLeft := float32(boxWidth * i)
			gtx.Constraints = layout.Exact(image.Pt(boxWidth, int(boxHeight)))
			fmt.Println(i, boxHeight, offsetTop, offsetLeft)
			trans := op.Offset(f32.Pt(offsetLeft, offsetTop)).Push(gtx.Ops)

			bar := layout.Inset{
				Left:  unit.Dp(b.SpaceBetween / 2),
				Right: unit.Dp(b.SpaceBetween / 2),
			}

			bar.Layout(gtx, fillWithLabel(b.Theme, b.BoxLabel(i), b.BarColor(i)))

			trans.Pop()
		}

		return layout.Dimensions{Size: gtx.Constraints.Max}
	}
}
