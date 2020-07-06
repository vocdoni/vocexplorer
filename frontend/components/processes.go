package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// ProcessView renders the processes page
type ProcessView struct {
	vecty.Core
}

// Render renders the ProcessView component
func (t *ProcessView) Render() vecty.ComponentOrHTML {
	return elem.Div(
		&Header{currentPage: "processes"},
	)
}
