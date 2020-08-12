package bootstrap

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/rpc"
)

type Badge struct {
	vecty.Core
	T    *rpc.TendermintInfo
	Text string
	Type string
}

func (b *Badge) Render() vecty.ComponentOrHTML {
	t := b.Type
	if len(b.Type) == 0 {
		t = "primary"
	}
	return elem.Span(
		vecty.Markup(vecty.Class("badge", "badge-"+t)),
		vecty.Text(
			b.Text,
		),
	)
}
