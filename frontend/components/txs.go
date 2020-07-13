package components

import (
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// TxsView renders the transaction page
type TxsView struct {
	vecty.Core
}

// Render renders the TxsView component
func (t *TxsView) Render() vecty.ComponentOrHTML {
	js.Global().Set("page", "transactions")
	js.Global().Set("gateway", false)

	return elem.Div(
		&Header{},
	)
}
