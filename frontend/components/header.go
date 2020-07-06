package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// Header renders the header
type Header struct {
	vecty.Core
	currentPage string
}

// Render renders the Header component
func (h *Header) Render() vecty.ComponentOrHTML {
	return elem.Div(
		elem.Heading1(
			vecty.Text("Vochain Block Explorer: "),
			vecty.Text(h.currentPage),
		),
		&NavBar{},
	)
}
