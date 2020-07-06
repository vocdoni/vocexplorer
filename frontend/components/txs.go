package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

// TxsView renders the transaction page
type TxsView struct {
	vecty.Core
}

// Render renders the TxsView component
func (t *TxsView) Render() vecty.ComponentOrHTML {
	return elem.Div(
		&Header{currentPage: "transactions"},
	)
}
