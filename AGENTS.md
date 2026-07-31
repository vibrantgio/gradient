# AGENTS.md — gradient

A two-stop linear gradient over the whole clip area: `FillLinearGradient`
paints into an op list, and `LinearGradient` returns a `layout.Widget`
whose two stops are given as fractions of the widget's own size rather than
in pixels.

**Layer.** Tier 0 of ADR-001's table — a leaf whose only dependency is Gio.
`mvu/example` is its only consumer anywhere in the organization.

**Read the canonical guide before you write code against this module.** It is
the organization's one agent guide — the module inventory with current tags,
the application skeleton, the MVU loop and rx semantics, typography, and the
pitfalls that are not guessable. It lives exactly once, in `vibrantgio/.github`,
and this file links it rather than copying it:

    https://raw.githubusercontent.com/vibrantgio/.github/master/llms.txt

**Module.** `github.com/vibrantgio/gradient`, one module at the repository
root.

**Build and test.** From the repository root:

    go build ./... && go test ./...
