package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// Header renders the header
type Header struct {
	vecty.Core
}

// Render renders the Header component
func (h *Header) Render() vecty.ComponentOrHTML {
	return elem.Header(
		&NavBar{},
	)
}
