# gradient

A two-stop linear gradient fill for [Vibrant Gio](https://github.com/vibrantgio),
a design system for native desktop applications on macOS, Windows and Linux,
written in pure Go on [Gio](https://gioui.org). Two functions, no state.

Gio has `paint.LinearGradientOp`, and it is not quite usable on its own. It
takes its two stops as absolute `f32.Point`s in the current coordinate space,
it paints nothing until you follow it with a `paint.PaintOp`, and it paints
into whatever clip happens to be current — so filling a component with a gradient
that runs corner to corner means knowing the component's pixel size, building the
clip yourself, and remembering the second op. That is four lines of ceremony
every time, and the pixel arithmetic has to be redone on every resize.

`LinearGradient` takes its stops as *fractions of the component's own size* —
`f32.Pt(0, 0)` to `f32.Pt(1, 1)` is corner to corner, whatever the window is
doing — clips to the incoming constraints, and hands back a `layout.Widget`.
`FillLinearGradient` is the same paint without the component, for callers that
already own a clip.

## Where it sits

Tier 0 of the stack — `mvu → theme → components → effects → patterns → markdown` —
a leaf whose only dependency is Gio. The
[organization page](https://github.com/vibrantgio) has the full tier table.

It sits alongside [backdrop](https://github.com/vibrantgio/backdrop) and
[circle](https://github.com/vibrantgio/circle), the other two drawing leaves.

```sh
go get github.com/vibrantgio/gradient
```

The Gio and Go versions are the ones `go.mod` names.

## Packages

One package, at the module root.

| Symbol | |
| --- | --- |
| `LinearGradient(stop1, color1, stop2, color2)` | A `layout.Widget` that clips to `gtx.Constraints.Max` and fills it. Stops are fractions of that size: `(0,0)` is the top-left corner, `(1,1)` the bottom-right. |
| `FillLinearGradient(ops, stop1, color1, stop2, color2)` | The same paint straight into an `*op.Ops`, with stops in absolute pixels and no clip of its own — it fills whatever clip is current. |

Both take `color.Color`, not `color.NRGBA`, and convert through
`color.NRGBAModel`, so a `colornames` value or a `color.RGBA` goes in
directly.

## Usage

A whole program — a window that is nothing but a gradient:

```go
window := mvu.NewWindow(app.Title("MVU - Gradient"))
gradient := rx.Of(gradient.LinearGradient(
	f32.Pt(0, 0), colornames.DeepPurple800,
	f32.Pt(1, 1), colornames.DeepPurple300,
))
window.Render(gradient).Wait()
```

`rx.Of` is what makes a static component into a layer: an MVU window renders
`rx.Observable[layout.Widget]` streams stacked back to front, and a background
that never changes is a stream of exactly one element. Add layers by passing
more observables to `Render`.

To paint a gradient inside a shape you already clipped — a rounded card, a
circle — drop to `FillLinearGradient`, whose stops are pixels rather than
fractions:

```go
defer clip.UniformRRect(bounds, radius).Push(gtx.Ops).Pop()
gradient.FillLinearGradient(gtx.Ops,
	f32.Pt(0, 0), top,
	f32.Pt(0, float32(bounds.Dy())), bottom,
)
```

## For coding assistants

Read this repository's [AGENTS.md](./AGENTS.md) before you write code against this module.

Build and test the root module with `go build ./...` and `go test ./...`.

## Status

Honest about what does not work yet.

- **Compositing many gradients into one shape wants a different shape.** A
  shadow drawn as a set of gradient tiles goes straight to
  `paint.LinearGradientOp`, because a component whose stops are fractions of
  the space it is handed cannot place several gradients inside one clip.
- **Two stops, linear only.** There is no multi-stop gradient, no angle
  parameter beyond choosing the two points, and no radial gradient: Gio
  exposes no radial primitive at all, so a radial look is approximated out of
  eight linear passes.
- **`LinearGradient` always takes all the space it is offered.** It clips to
  and returns `gtx.Constraints.Max`, ignoring `Constraints.Min`, so as a flex
  child it fills the flex rather than sizing to anything. It is a background
  layer, not a component.
- **`FillLinearGradient` paints the current clip, whatever that is.** Called
  with no clip pushed, it fills the entire window. That is the same contract as
  [backdrop](https://github.com/vibrantgio/backdrop)'s `Fill`, and it is easy to
  use by accident when the component form is what you wanted.
- **Colours are arguments, not colour roles.** Neither function knows about
  `tokens.PlatformColors`, so a gradient does not follow the theme unless the
  caller re-derives its two colours on every theme emission and rebuilds the
  component. This module has no token-aware wrapper.
- **There are no tests and no golden images.** `go test ./...` reports "no test
  files", so nothing pins the fraction-to-pixel arithmetic.

## License

MIT — see [LICENSE](./LICENSE).
