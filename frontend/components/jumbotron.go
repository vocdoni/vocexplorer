package components

import (
	"fmt"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"gitlab.com/vocdoni/vocexplorer/frontend/store"
	"gitlab.com/vocdoni/vocexplorer/util"
)

//Jumbotron is a component for an info banner of statistics
type Jumbotron struct {
	vecty.Core
}

//JumboStatTitle renders the jumbotron statistic title
func JumboStatTitle(t string) vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(vecty.Class("stat-title")),
		vecty.Text(t),
	)
}

//JumboStatValue renders a jumbotron statistic contents
func JumboStatValue(v string) vecty.MarkupOrChild {
	return elem.Div(
		vecty.Markup(vecty.Class("stat-value")),
		vecty.Text(v),
	)
}

//Render renders the jumbotron component
func (b *Jumbotron) Render() vecty.ComponentOrHTML {
	colMarkup := vecty.Markup(vecty.Class("col-xs-12", "col-sm-4", "mb-2", "mb-sm-0"))
	var items vecty.List

	if store.Stats.BlockTime != nil && store.Stats.BlockTime[0] > 0 {
		items = append(items, elem.Div(
			colMarkup,
			JumboStatTitle("Average block time"),
			JumboStatValue(fmt.Sprintf("%ss", util.MsToSecondsString(store.Stats.BlockTime[0]))),
		))
	}

	if store.Processes.Count > 0 {
		items = append(items, elem.Div(
			colMarkup,
			JumboStatTitle("Total processes"),
			JumboStatValue(util.IntToString(store.Processes.Count)),
		))
	}

	if store.Entities.Count > 0 {
		items = append(items, elem.Div(
			colMarkup,
			JumboStatTitle("Total entities"),
			JumboStatValue(util.IntToString(store.Entities.Count)),
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
