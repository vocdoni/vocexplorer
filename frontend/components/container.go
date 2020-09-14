package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

// Container defines div with a .container class
func Container(markup ...vecty.MarkupOrChild) vecty.ComponentOrHTML {
	markup = append(
		markup,
		vecty.Markup(vecty.Class("container")),
	)
	return elem.Div(markup...)
}
