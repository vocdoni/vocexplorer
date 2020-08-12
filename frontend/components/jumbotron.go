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
	var items []vecty.ComponentOrHTML

	if b.vc.BlockTime != nil && b.vc.BlockTime[0] > 0 {
		items = append(items, elem.Div(
			colMarkup,
			JumboStatTitle("Average block time"),
			JumboStatValue(fmt.Sprintf("%ss", util.MsToSecondsString(b.vc.BlockTime[0]))),
		))
	}

	// to be replaced with its proper condition whenever we have the process count
	if true {
		items = append(items, elem.Div(
			colMarkup,
			JumboStatTitle("Total processes"),
			JumboStatValue("23,418"),
		))
	}

	// to be replaced with its proper condition whenever we have the process count
	if true {
		items = append(items, elem.Div(
			colMarkup,
			JumboStatTitle("Total entities"),
			JumboStatValue("2,323"),
		))
	}

	return elem.Div(
		vecty.Markup(vecty.Class("jumbotron")),
		elem.Div(
			vecty.Markup(vecty.Class("container")),
			elem.Div(
				vecty.Markup(vecty.Class("jumbo-stats", "row")),
				vecty.List(items),
			),
		),
	)
}
