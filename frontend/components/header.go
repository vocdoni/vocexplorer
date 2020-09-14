package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
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
