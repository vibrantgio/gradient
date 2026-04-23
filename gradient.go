package gradient

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func FillLinearGradient(ops *op.Ops, stop1 f32.Point, color1 color.Color, stop2 f32.Point, color2 color.Color) {
	paint.LinearGradientOp{
		Stop1:  stop1,
		Color1: color.NRGBAModel.Convert(color1).(color.NRGBA),
		Stop2:  stop2,
		Color2: color.NRGBAModel.Convert(color2).(color.NRGBA),
	}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func LinearGradient(stop1 f32.Point, color1 color.Color, stop2 f32.Point, color2 color.Color) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		size := gtx.Constraints.Max
		defer clip.Outline{Path: clip.Rect{Max: size}.Path()}.Op().Push(gtx.Ops).Pop()
		FillLinearGradient(gtx.Ops, f32.Pt(stop1.X*float32(size.X), stop1.Y*float32(size.Y)), color1, f32.Pt(stop2.X*float32(size.X), stop2.Y*float32(size.Y)), color2)
		return layout.Dimensions{Size: size}
	}
}
