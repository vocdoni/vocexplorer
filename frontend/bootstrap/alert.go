package bootstrap

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

type Alert struct {
	vecty.Core
	Contents string
	Type     string
}

func (b *Alert) Render() vecty.ComponentOrHTML {
	t := b.Type
	if len(b.Type) == 0 {
		t = "success"
	}

	return elem.Div(
		vecty.Markup(vecty.Class("alert", "alert-"+t)),
		vecty.Markup(
			vecty.UnsafeHTML(b.Contents),
		),
	)
}
