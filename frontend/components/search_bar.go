package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

//SearchBar is a component for a search bar
type SearchBar struct {
	vecty.Core
}

//Render renders the SearchBar component
func (s *SearchBar) Render() vecty.ComponentOrHTML {
	return elem.Form(
		vecty.Markup(vecty.Class("form-inline", "my-2", "my-lg-0")),
		elem.Div(
			vecty.Markup(vecty.Class("input-group")),
			elem.Div(
				vecty.Markup(vecty.Class("input-group-append")),
				elem.Button(
					vecty.Markup(vecty.Class("btn", "input-group-text")),
					vecty.Markup(vecty.Attribute("type", "submit")),
					elem.Span(
						vecty.Markup(vecty.Class("sr-only")),
						vecty.Text("Search"),
					),
					elem.Span(
						vecty.Markup(vecty.Class("icon-lens")),
					),
				),
			),
			elem.Input(
				vecty.Markup(vecty.Class("form-control", "mr-sm-2")),
				vecty.Markup(vecty.Attribute("aria-label", "Table of average block times for given time periods.")),
				vecty.Markup(vecty.Attribute("type", "search")),
				vecty.Markup(vecty.Attribute("placeholder", "Search by process, entity, transaction, block height or address")),
			),
		),
	)
}
