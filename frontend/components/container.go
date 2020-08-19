package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// Container defines div with a .container class
func Container(markup ...vecty.MarkupOrChild) vecty.ComponentOrHTML {
	markup = append(
		markup,
		vecty.Markup(vecty.Class("container")),
	)
	return elem.Div(markup...)
}
