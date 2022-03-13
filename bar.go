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
	"gioui.org/text"
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

	return b.renderYLegend(gtx, b.Theme)
}

func (b Bar) renderYLegend(gtx layout.Context, th *material.Theme) layout.Dimensions {
	_, maxValue := minMax(b.Data)

	lineWidth := 10
	leftSize := 100
	rightSize := gtx.Constraints.Max.X - leftSize

	var barsDims layout.Dimensions
	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(rightSize, gtx.Constraints.Max.Y))
		trans := op.Offset(f32.Pt(float32(leftSize), 0)).Push(gtx.Ops)
		barsDims = layout.Inset{
			Top:    th.TextSize,
			Bottom: th.TextSize,
		}.Layout(gtx, b.renderBars(maxValue))
		trans.Pop()
	}
	{
		gtx := gtx
		gtx.Constraints = layout.Exact(image.Pt(leftSize, gtx.Constraints.Max.Y))

		writeYLegendLine(gtx, lineWidth)
		gtx.Constraints.Max.X -= lineWidth + 2

		labelFunc := func(v float64) layout.Widget {
			l := material.Label(th, th.TextSize, fmt.Sprintf("%.2f", v))
			l.Alignment = text.End
			l.MaxLines = 1
			return l.Layout
		}

		maxBarHeight := barsDims.Size.Y - 2*barsDims.Baseline

		actualHeight := gtx.Constraints.Min.Y - 48
		lineHeight := int(th.TextSize.V * 1.5)

		if actualHeight > lineHeight*2 {
			layout.Inset{}.Layout(gtx, labelFunc(maxValue))
		}
		if actualHeight > lineHeight*10 {
			layout.Inset{
				Top: unit.Px(float32(maxBarHeight / 4)),
			}.Layout(gtx, labelFunc(maxValue/4*3))
		}
		if actualHeight > lineHeight*6 {
			layout.Inset{
				Top: unit.Px(float32(maxBarHeight / 2)),
			}.Layout(gtx, labelFunc(maxValue/2))
		}
		if actualHeight > lineHeight*10 {
			layout.Inset{
				Top: unit.Px(float32(maxBarHeight / 4 * 3)),
			}.Layout(gtx, labelFunc(maxValue/4))
		}
		if actualHeight > lineHeight*2 {
			layout.Inset{
				Top: unit.Px(float32(maxBarHeight)),
			}.Layout(gtx, labelFunc(0))
		}
	}

	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func writeYLegendLine(gtx layout.Context, width int) layout.Dimensions {
	size := image.Pt(2, gtx.Constraints.Max.Y)
	defer clip.Rect{
		Min: image.Pt(gtx.Constraints.Max.X-width, 0),
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

		spaceBetween := unit.Dp(b.SpaceBetween / 2)
		if gtx.Constraints.Max.X < len(b.Data)*2*int(spaceBetween.V)*5 {
			spaceBetween = unit.Dp(0)
		}
		for i, n := range b.Data {
			gtx := gtx
			boxHeight := (n / maxN) * float64(gtx.Constraints.Max.Y)
			offsetTop := float32(gtx.Constraints.Max.Y) - float32(boxHeight)
			offsetLeft := float32(boxWidth * i)
			gtx.Constraints = layout.Exact(image.Pt(boxWidth, int(boxHeight)))
			trans := op.Offset(f32.Pt(offsetLeft, offsetTop)).Push(gtx.Ops)

			bar := layout.Inset{
				Left:  spaceBetween,
				Right: spaceBetween,
			}

			bar.Layout(gtx, fillWithLabel(b.Theme, b.BoxLabel(i), b.BarColor(i)))

			trans.Pop()
		}

		return layout.Dimensions{Size: gtx.Constraints.Max}
	}
}
