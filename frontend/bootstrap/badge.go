package bootstrap

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type Badge struct {
	vecty.Core
	Contents string
	Type     string
}

func (b *Badge) Render() vecty.ComponentOrHTML {
	t := b.Type
	if len(b.Type) == 0 {
		t = "primary"
	}
	return elem.Span(
		vecty.Markup(vecty.Class("badge", "badge-"+t)),
		vecty.Markup(
			vecty.UnsafeHTML(b.Contents),
		),
	)
}
