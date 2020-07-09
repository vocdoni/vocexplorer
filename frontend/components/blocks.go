package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// BlocksView renders the blocks page
type BlocksView struct {
	vecty.Core
}

// Render renders the BlocksView component
func (t *BlocksView) Render() vecty.ComponentOrHTML {
	js.Global().Set("page", "blocks")
	js.Global().Set("gateway", false)
	return elem.Div(
		&Header{currentPage: "blocks"},
	)
}
