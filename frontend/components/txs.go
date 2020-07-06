package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/nathanhack/vectyUI/nav"
)

// TxsView renders the transation page
type TxsView struct {
	vecty.Core
}

// Render renders the TxsView component
func (t *TxsView) Render() vecty.ComponentOrHTML {
	return elem.Div{
		NavBar{}
	}
}

