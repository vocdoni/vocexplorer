package components

import (
	"fmt"

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

func JumboStatValue(v string) vecty.MarkupOrChild {
	return elem.Div(
		vecty.Markup(vecty.Class("stat-value")),
		vecty.Text(v),
	)
}

func (b *Jumbotron) Render() vecty.ComponentOrHTML {
	colMarkup := vecty.Markup(vecty.Class("col-xs-12", "col-sm-4"))
	var items vecty.List

	if b.vc.BlockTime != nil && b.vc.BlockTime[0] > 0 {
		items = append(items, elem.Div(
			colMarkup,
			JumboStatTitle("Average block time"),
			JumboStatValue(fmt.Sprintf("%ss", util.MsToSecondsString(b.vc.BlockTime[0]))),
		))
	}

	if b.vc.ProcessCount > 0 {
		items = append(items, elem.Div(
			colMarkup,
			JumboStatTitle("Total processes"),
			JumboStatValue(util.IntToString(b.vc.ProcessCount)),
		))
	}

	if b.vc.EntityCount > 0 {
		items = append(items, elem.Div(
			colMarkup,
			JumboStatTitle("Total entities"),
			JumboStatValue(util.IntToString(b.vc.EntityCount)),
		))
	}

	return elem.Div(
		vecty.Markup(vecty.Class("jumbotron")),
		Container(
			elem.Div(
				vecty.Markup(vecty.Class("jumbo-stats", "row")),
				vecty.List(items),
			),
		),
	)
}
