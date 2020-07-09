package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// ProcessView renders the processes page
type ProcessView struct {
	vecty.Core
}

// Render renders the ProcessView component
func (t *ProcessView) Render() vecty.ComponentOrHTML {
	js.Global().Set("page", "processes")
	js.Global().Set("gateway", false)

	return elem.Div(
		&Header{currentPage: "processes"},
	)
}
