package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/client"
	"gitlab.com/vocdoni/vocexplorer/rpc"
	"gitlab.com/vocdoni/vocexplorer/util"
)

type Jumbotron struct {
	vecty.Core
	t  *rpc.TendermintInfo
	vc *client.VochainInfo
}

func JumboStatTitle(t string) vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("stat-title")),
		vecty.Text(t),
	)
}

func JumboStatValue(c bool, t string) vecty.MarkupOrChild {
	return vecty.If(
		c,
		elem.Div(
			vecty.Markup(vecty.Class("stat-value")),
			vecty.Text(t),
		),
	)
}

func (b *Jumbotron) Render() vecty.ComponentOrHTML {
	colMarkup := vecty.Markup(vecty.Class("col-sm-12", "col-lg-6", "col-xl-3"))
	return elem.Div(
		vecty.Markup(vecty.Class("jumbotron")),
		elem.Div(
			vecty.Markup(vecty.Class("container")),
			elem.Div(
				vecty.Markup(vecty.Class("jumbo-stats", "row")),
				elem.Div(
					colMarkup,
					JumboStatTitle("Average block time"),
					JumboStatValue(b.vc.BlockTime != nil && b.vc.BlockTime[0] > 0 && b.t.ResultStatus != nil, "5s"),
				),
				elem.Div(
					colMarkup,
					JumboStatTitle("Total processes"),
					JumboStatValue(true, util.IntToString(b.vc.ProcessCount)),
				),
				elem.Div(
					colMarkup,
					JumboStatTitle("Total entities"),
					JumboStatValue(true, util.IntToString(b.vc.EntityCount)),
				),
			),
		),
	)
}
