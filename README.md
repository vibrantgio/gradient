# gradient

A two-stop linear gradient fill for [Vibrant Gio](https://github.com/vibrantgio),
a design system for native desktop applications on macOS, Windows and Linux,
written in pure Go on [Gio](https://gioui.org). Two functions, thirty lines,
no state.

Gio has `paint.LinearGradientOp`, and it is not quite usable on its own. It
takes its two stops as absolute `f32.Point`s in the current coordinate space,
it paints nothing until you follow it with a `paint.PaintOp`, and it paints
into whatever clip happens to be current — so filling a widget with a gradient
that runs corner to corner means knowing the widget's pixel size, building the
clip yourself, and remembering the second op. That is four lines of ceremony
every time, and the pixel arithmetic has to be redone on every resize.

`LinearGradient` takes its stops as *fractions of the widget's own size* —
`f32.Pt(0, 0)` to `f32.Pt(1, 1)` is corner to corner, whatever the window is
doing — clips to the incoming constraints, and hands back a `layout.Widget`.
`FillLinearGradient` is the same paint without the widget, for callers that
already own a clip.

## Where it sits

Tier 0 of the stack — `mvu → spectrum → prism → pulse → cadence → markdown` —
a leaf whose only dependency is Gio. The
[organization page](https://github.com/vibrantgio) has the full tier table.

It sits alongside [backdrop](https://github.com/vibrantgio/backdrop) and
[circle](https://github.com/vibrantgio/circle), the other two drawing leaves:
the three of them are what an MVU application composes its bottom layer out of
before any component is involved.

```sh
go get github.com/vibrantgio/gradient
```

Every module in the organization is on gioui.org v0.10.1 and Go 1.25.1.

## Packages

One package, at the module root.

| Symbol | |
| --- | --- |
| `LinearGradient(stop1, color1, stop2, color2)` | A `layout.Widget` that clips to `gtx.Constraints.Max` and fills it. Stops are fractions of that size: `(0,0)` is the top-left corner, `(1,1)` the bottom-right. |
| `FillLinearGradient(ops, stop1, color1, stop2, color2)` | The same paint straight into an `*op.Ops`, with stops in absolute pixels and no clip of its own — it fills whatever clip is current. |

Both take `color.Color`, not `color.NRGBA`, and convert through
`color.NRGBAModel`, so a `colornames` value or an `color.RGBA` goes in
directly.

## Usage

The whole of `mvu/example/03-gradient` — a window that is nothing but a
gradient:

```go
window := mvu.NewWindow(app.Title("MVU - Gradient"))
gradient := rx.Of(gradient.LinearGradient(
	f32.Pt(0, 0), colornames.DeepPurple800,
	f32.Pt(1, 1), colornames.DeepPurple300,
))
window.Render(gradient).Wait()
```

`rx.Of` is what makes a static widget into a layer: an MVU window renders
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

Read the canonical guide before writing code against this module — the module
inventory with current tags, the application skeleton, MVU and rx semantics,
typography, and the pitfalls that are not guessable:

<https://raw.githubusercontent.com/vibrantgio/.github/master/llms.txt>

[`AGENTS.md`](./AGENTS.md) in this repository has the build and test commands.

## Status

Honest about what does not work yet. Every count below is measured.

- **One consumer, one call site.** `mvu/example/03-gradient` is the only
  program anywhere in the organization that imports this module, and
  `LinearGradient` is the only symbol it uses. `FillLinearGradient` has never
  been called outside this repository.
- **The effects layer needed gradients and did not use this one.** `pulse/depth`
  composes a cast shadow out of eight linear gradients and `pulse/glow` composes
  a halo out of eight more, and both call `paint.LinearGradientOp` directly
  rather than importing this module — the fractional-stop widget is the wrong
  shape for compositing many gradients into one shape. Nothing in the current
  plan reconciles the two.
- **Two stops, linear only.** There is no multi-stop gradient, no angle
  parameter beyond choosing the two points, and no radial gradient — Gio
  exposes no radial primitive at all, which is why `pulse` fakes one out of
  eight linear passes. Phase E of the
  [org plan](https://github.com/vibrantgio/.github) builds a blur on
  `gioui.org/gpu/headless` and revisits that; it does not claim this module.
- **`LinearGradient` always takes all the space it is offered.** It clips to
  and returns `gtx.Constraints.Max`, ignoring `Constraints.Min`, so as a flex
  child it fills the flex rather than sizing to anything. It is a background
  layer, not a component.
- **`FillLinearGradient` paints the current clip, whatever that is.** Called
  with no clip pushed, it fills the entire window. That is the same contract as
  [backdrop](https://github.com/vibrantgio/backdrop)'s `Fill`, and it is easy to
  reach for by accident when the widget form is what you wanted.
- **Colours are arguments, not roles.** Neither function knows about
  `tokens.ColorTokens`, so a gradient does not follow the theme unless the
  caller re-derives its two colours on every theme emission and rebuilds the
  widget. There is no token-aware wrapper anywhere in the organization.
- **There are no tests and no golden images.** `go test ./...` reports "no test
  files", so nothing pins the fraction-to-pixel arithmetic.
- **There is no LICENSE file.** Eighteen of the organization's twenty
  repositories ship one — MIT for most, BSD, Apache or Unlicense for the ported
  libraries. This one and [circle](https://github.com/vibrantgio/circle) are
  the two that do not, and no phase of the current plan adds them.
